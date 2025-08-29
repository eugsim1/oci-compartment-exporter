
````markdown
# OCI Compartment Exporter

This Go program lists all compartments in an Oracle Cloud Infrastructure (OCI) tenancy and generates a CSV file with the full paths, levels, and parent-child relationships. It can optionally filter for a specific subtree based on a root compartment OCID.

## Features

- Lists all compartments in the tenancy, including nested compartments.
- Builds full paths for each compartment (e.g., `ROOT/Dept/Team`).
- Calculates the level of each compartment in the hierarchy.
- Optionally filters compartments starting from a specific root compartment.
- Outputs results in a CSV file `oci_compartment_paths.csv`.

## Prerequisites

- Go 1.18+ installed
- Oracle OCI Go SDK v65
- OCI CLI config file (`~/.oci/config`) or environment variables set for authentication
- Access to the OCI tenancy

## Installation

1. Clone the repository:

```bash
git clone https://github.com/yourusername/oci-compartment-exporter.git
cd oci-compartment-exporter
````

2. Initialize Go modules:

```bash
go mod tidy
```

3. Build the project:

```bash
go build -o oci-compartment-exporter main.go
```

## Usage

```bash
./oci-compartment-exporter [-root <root_compartment_ocid>]
```

### Flags

* `-root` (optional): OCID of a compartment to limit the output to its subtree. If not provided, all compartments in the tenancy will be listed.

### Example

Export all compartments in the tenancy:

```bash
./oci-compartment-exporter
```

Export compartments under a specific root compartment:

```bash
./oci-compartment-exporter -root ocid1.compartment.oc1..exampleuniqueID
```

## Output

The program generates a CSV file `oci_compartment_paths.csv` with the following columns:

* `id`: Compartment OCID
* `parent_id`: Parent compartment OCID
* `level`: Depth level in the hierarchy (root = 0)
* `path`: Full path from the root compartment

Example:

| id                                     | parent\_id                            | level | path                 |
| -------------------------------------- | ------------------------------------- | ----- | -------------------- |
| ocid1.tenancy.oc1..exampleuniqueID     |                                       | 0     | ROOT                 |
| ocid1.compartment.oc1..exampleChildID  | ocid1.tenancy.oc1..exampleuniqueID    | 1     | ROOT/Finance         |
| ocid1.compartment.oc1..exampleChildID2 | ocid1.compartment.oc1..exampleChildID | 2     | ROOT/Finance/Payroll |

## Notes

* Ensure your OCI user has the necessary permissions to list compartments.
* The CSV file will be created in the same directory where the program is executed.
* Pagination is handled automatically.

## License

This project is licensed under the MIT License.

```


