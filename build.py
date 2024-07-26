import requests
import os
from dotenv import load_dotenv

# Load environment variables from .env file
load_dotenv()

# Inputs
release_name = input("Release Name: ").strip()
release_body = input("Release Body: ").strip()
tag_name = input("Tag Name: ").strip()

# Check if the inputs were filled
if not release_name or not release_body or not tag_name:
    print("Please fill in all required inputs.")
    exit(1)

# Get GitHub token from environment variables
token = os.getenv("GITHUB_TOKEN")
if not token:
    print("GitHub token not found. Please set GITHUB_TOKEN in your .env file.")
    exit(1)

repo_name = "Purrquinox/Chlamydia"
bin_directory = "bin"

# Create a release
release_url = f"https://api.github.com/repos/{repo_name}/releases"
headers = {
    "Authorization": f"Bearer {token}",
    "Accept": "application/vnd.github.v3+json",
    "X-GitHub-Api-Version": "2022-11-28"
}
data = {
    "tag_name": tag_name,
    "name": release_name,
    "body": release_body,
    "draft": False,
    "prerelease": False,
    "latest": True
}

# Send the request to create the release
response = requests.post(release_url, headers=headers, json=data)

# Check the response status code and handle errors
if response.status_code == 201:
    release = response.json()
    print(f"Release created: {release['html_url']}")

    upload_url = release.get("upload_url", "").replace("{?name,label}", "")
    if not upload_url:
        print("Upload URL not found in the response.")
    else:
        # Recursively find and list all files in the bin directory
        file_paths = []
        for root, dirs, files in os.walk(bin_directory):
            for file in files:
                file_paths.append(os.path.join(root, file))

        if not file_paths:
            print(f"No files found in directory {bin_directory}.")
            exit(1)

        print(f"Found files: {file_paths}")

        # Upload files to the release
        upload_headers = {
            "Authorization": f"Bearer {token}",
            "Content-Type": "application/octet-stream",
        }

        for file_path in file_paths:
            # Create a file name with subdirectory type and file name
            relative_path = os.path.relpath(file_path, bin_directory)
            dir_name, file_name = os.path.split(relative_path)
            new_file_name = f"{dir_name}-{file_name}"  # Include subdirectory in the file name
            
            with open(file_path, "rb") as file:
                upload_response = requests.post(
                    f"{upload_url}?name={new_file_name}",
                    headers=upload_headers,
                    data=file
                )
            # Check the response for each file upload
            if upload_response.status_code == 201:
                print(f"Successfully uploaded {new_file_name}")
            else:
                print(f"Failed to upload {new_file_name}: {upload_response.status_code}")
                print(upload_response.json())
else:
    print(f"Failed to create release: {response.status_code}")
    print(response.json())
