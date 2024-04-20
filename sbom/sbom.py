import json
import os

# Define the servers and their respective SBOM file names
servers = {
    "user": "user-server-sbom.json",
    "post": "post-server-sbom.json",
    "match": "match-server-sbom.json",
    "bid": "bid-server-sbom.json",
    "verification": "verification-server-sbom.json"
}

# Function to extract licenses from a component
def extract_licenses(component):
    licenses = []
    if "evidence" in component and "licenses" in component["evidence"]:
        for license_entry in component["evidence"]["licenses"]:
            if "license" in license_entry:
                license_id = license_entry["license"]["id"]
                licenses.append(license_id)
    return licenses

# Function to generate licenses list for a server
def generate_licenses_list(server_name):
    sbom_file = os.path.join(f"{server_name}-server-sbom.json")

    try:
        with open(sbom_file, "r") as f:
            sbom_data = json.load(f)
    except FileNotFoundError:
        print(f"SBOM file not found: {sbom_file}")
        return

    licenses = set()
    for component in sbom_data["components"]:
        component_licenses = extract_licenses(component)
        licenses.update(component_licenses)

    print(f"{server_name}-server licenses:")
    for license_id in sorted(licenses):
        print(f"- {license_id}")
    print()

# Generate licenses lists for all servers
for server_name, sbom_file in servers.items():
    generate_licenses_list(server_name)
