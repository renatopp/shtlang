import yaml
import json
import sys

def convert_tmLanguage_yaml_to_json(yaml_file, json_file):
    try:
        # Load data from the YAML file
        with open(yaml_file, 'r') as file:
            data = yaml.safe_load(file)

        # Write the data to a JSON file
        with open(json_file, 'w') as file:
            json.dump(data, file, indent=2)

        print(f"Conversion successful. Data saved to {json_file}")
    except Exception as e:
        print(f"An error occurred: {e}")

if __name__ == "__main__":
    # Specify the file paths
    input_yaml = 'src/sht.tmLanguage.yaml'
    output_json = 'syntaxes/sht.tmLanguage.json'

    # Convert the YAML file to JSON
    convert_tmLanguage_yaml_to_json(input_yaml, output_json)