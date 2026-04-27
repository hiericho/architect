import os
import subprocess
import sys
import shutil

# Target platforms: (GOOS, GOARCH, extension)
TARGETS = [
    ("windows", "amd64", ".exe"),
    ("linux", "amd64", ""),
    ("linux", "arm64", ""),
    ("darwin", "amd64", ""), # macOS Intel
    ("darwin", "arm64", ""), # macOS Apple Silicon
]

def main():
    script_dir = os.path.dirname(os.path.abspath(__file__))
    root_dir = os.path.dirname(script_dir)
    engine_dir = os.path.join(root_dir, "engine")
    bin_dir = os.path.join(root_dir, "architect", "bin")

    # Ensure bin directory exists
    if not os.path.exists(bin_dir):
        os.makedirs(bin_dir)

    print("🏗️ Building Architect Engine for all platforms...")

    for goos, goarch, ext in TARGETS:
        bin_name = f"architect_{goos}_{goarch}{ext}"
        bin_path = os.path.join(bin_dir, bin_name)
        
        print(f"  -> Compiling {bin_name}...")
        
        env = os.environ.copy()
        env["GOOS"] = goos
        env["GOARCH"] = goarch
        env["CGO_ENABLED"] = "0" # Ensure static linking

        try:
            subprocess.run(
                ["go", "build", "-trimpath", "-ldflags", "-s -w", "-o", bin_path, "./main.go"],
                cwd=engine_dir,
                env=env,
                check=True
            )
            print(f"     ✅ Success")
        except subprocess.CalledProcessError as e:
            print(f"     ❌ Failed: {e}")
            sys.exit(1)

    print("\n🎉 All binaries built successfully in architect/bin/")
    print("📦 You can now run 'python setup.py sdist bdist_wheel' to package the project.")

if __name__ == "__main__":
    main()
