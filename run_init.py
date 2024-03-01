import os
import subprocess
import webbrowser

def is_docker_installed():
    try:
        subprocess.run(["docker", "--version"], check=True)
        return True
    except Exception:
        return False

def install_docker():
    print("Docker is not installed. Please install Docker first.")
    # Additional instructions or automation for Docker installation can go here.

def create_dockerfile(create_static_files):
    dockerfile_lines = ["FROM nginx:latest"]
    
    # Check if nginx.conf exists before adding it to Dockerfile.
    if os.path.exists("nginx.conf"):
        dockerfile_lines.append("COPY nginx.conf /etc/nginx/nginx.conf")
    
    if create_static_files:
        dockerfile_lines.append("COPY ./html /usr/share/nginx/html")
    
    dockerfile_lines += ["EXPOSE 80", 'CMD ["nginx", "-g", "daemon off;"]']
    
    with open("Dockerfile", "w") as dockerfile:
        dockerfile.write("\n".join(dockerfile_lines))
    
    if create_static_files:
        os.makedirs("html", exist_ok=True)
        # Optionally, guide the user to populate this directory.

def build_docker_image():
    result = subprocess.run(["docker", "build", "-t", "cdn-nginx-image", "."], capture_output=True, text=True)
    return result.returncode == 0, result.stdout + result.stderr

def run_docker_container():
    result = subprocess.run(["docker", "run", "-d", "-p", "8080:80", "--name", "my-nginx-container", "cdn-nginx-image"], capture_output=True, text=True)
    return result.returncode == 0, result.stdout + result.stderr

def clear_screen():
    os.system('cls' if os.name == 'nt' else 'clear')

def open_localhost_in_browser():
    webbrowser.open("http://localhost:8080")

def main():
    if not is_docker_installed():
        install_docker()
        return

    create_static = input("Do you want to create a directory for static files? (y/n): ").lower() == 'y'
    create_dockerfile(create_static)
    
    print("Building Docker image...")
    success, build_logs = build_docker_image()
    if not success:
        print("Error building Docker image:\n", build_logs)
        return
    
    print("Running Docker container...")
    success, run_logs = run_docker_container()
    if not success:
        print("Error running Docker container:\n", run_logs)
        return
    
    clear_screen()
    if success:
        print("Successfully built and ran the Docker container.")
        open_localhost_in_browser()
    else:
        print("Failed to build or run the Docker container.")

if __name__ == "__main__":
    main()
