import os


def count_lines(file_path):
    with open(file_path, "r", encoding="utf-8") as file:
        lines = file.readlines()
    return len(lines)


def count_lines_in_folder(folder_path, forbidden_dirs, extensions, output_file):
    total_lines = 0
    with open(output_file, "w", encoding="utf-8") as output:
        for root, dirs, files in os.walk(folder_path):
            dirs[:] = [d for d in dirs if d not in forbidden_dirs]
            for file in files:
                file_path = os.path.join(root, file)
                if any(file.endswith(ext) for ext in extensions):
                    lines = count_lines(file_path)
                    total_lines += lines
                    with open(file_path, "r", encoding="utf-8") as input_file:
                        output.write(f"File: {file_path}\n")
                        for fline in input_file.readlines():
                            if fline.strip():
                                output.write(fline)
                        output.write("\n" + "=" * 50 + "\n")

    return total_lines


if __name__ == "__main__":
    start_path = input("start path (default: .): ")
    if start_path == "":
        start_path = "."

    forbidden_dirs = input("enter a forbidden dirs (etc: '.git, node_modules, bin'): ")
    forbidden_dirs = [d.strip() for d in forbidden_dirs.split(",")]

    extensions = input("enter list of extensions to include (etc: '.go, .ts, .cpp'): ")
    extensions = [d.strip() for d in extensions.split(",")]

    output_file = "output.txt"

    total_lines_of_code = count_lines_in_folder(
        start_path, forbidden_dirs, extensions, output_file
    )
    print(f"Total lines of code: {total_lines_of_code}")
    print(f"Code saved to: {output_file}")
