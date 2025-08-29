package main

import (
	"context"
	"encoding/csv"
	"flag"
	"fmt"
	"log"
	"os"
	"sort"
	"strings"

	"github.com/oracle/oci-go-sdk/v65/common"
	"github.com/oracle/oci-go-sdk/v65/identity"
)

func main() {
	rootCompartmentID := flag.String("root", "", "OCID of the root compartment to dump only its subtree (optional)")
	flag.Parse()

	// Load OCI config
	provider := common.DefaultConfigProvider()

	// Create Identity client
	client, err := identity.NewIdentityClientWithConfigurationProvider(provider)
	if err != nil {
		log.Fatalf("Error creating OCI identity client: %v", err)
	}

	compartmentMap := make(map[string]string) // ID -> Name
	parentMap := make(map[string]string)      // ID -> ParentID

	// Get root tenancy OCID
	
	tenancyID, err := provider.TenancyOCID()
	if err != nil {
		fmt.Errorf("failed to read tenancy OCID from config: %w", err)
	}
	

	compartmentMap[tenancyID] = "ROOT"
	parentMap[tenancyID] = "" // ROOT has no parent

	// List all compartments (handle pagination)
	request := identity.ListCompartmentsRequest{
		CompartmentId:          &tenancyID,
		CompartmentIdInSubtree: common.Bool(true),
		AccessLevel:            identity.ListCompartmentsAccessLevelAny,
	}

	page := ""
	for {
		request.Page = &page
		resp, err := client.ListCompartments(context.Background(), request)
		if err != nil {
			log.Fatalf("Error listing compartments: %v", err)
		}

		for _, c := range resp.Items {
			compartmentMap[*c.Id] = *c.Name
			parentMap[*c.Id] = *c.CompartmentId
		}

		if resp.OpcNextPage == nil || *resp.OpcNextPage == "" {
			break
		}
		page = *resp.OpcNextPage
	}

	// Build full paths
	type CompartmentPath struct {
		ID       string
		ParentID string
		Level    int
		Path     string
	}
	var paths []CompartmentPath

	for id := range compartmentMap {
		parentID := parentMap[id]
		path, level := buildFullPathAndLevel(id, compartmentMap, parentMap)

		// If rootCompartmentID is set, skip compartments outside this subtree
		if *rootCompartmentID != "" && !strings.HasPrefix(path, buildFullPath(*rootCompartmentID, compartmentMap, parentMap)) && id != *rootCompartmentID {
			continue
		}

		paths = append(paths, CompartmentPath{
			ID:       id,
			ParentID: parentID,
			Level:    level,
			Path:     path,
		})
	}

	// Sort by Level, then Path
	sort.Slice(paths, func(i, j int) bool {
		if paths[i].Level != paths[j].Level {
			return paths[i].Level < paths[j].Level
		}
		return paths[i].Path < paths[j].Path
	})

	// Write CSV
	file, err := os.Create("oci_compartment_paths.csv")
	if err != nil {
		log.Fatalf("Failed to create CSV file: %v", err)
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	// CSV header
	writer.Write([]string{"id", "parent_id", "level", "path"})

	for _, c := range paths {
		writer.Write([]string{c.ID, c.ParentID, fmt.Sprintf("%d", c.Level), c.Path})
	}

	fmt.Println("CSV file generated successfully: oci_compartment_paths.csv")
}

// Recursively build full path and level
func buildFullPathAndLevel(compartmentID string, names map[string]string, parents map[string]string) (string, int) {
	name, ok := names[compartmentID]
	if !ok {
		return "", 0
	}
	parentID, hasParent := parents[compartmentID]
	if !hasParent || parentID == compartmentID || parentID == "" {
		return name, 0
	}
	parentPath, level := buildFullPathAndLevel(parentID, names, parents)
	return parentPath + "/" + name, level + 1
}

// Build path only (for filtering)
func buildFullPath(compartmentID string, names map[string]string, parents map[string]string) string {
	path, _ := buildFullPathAndLevel(compartmentID, names, parents)
	return path
}
