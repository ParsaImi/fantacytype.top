fisrfile = input("Enter the first filename: ")
secondfile = input("Enter the second filename: ")

lines = set()
with open(fisrfile, "r") as f1, open(secondfile, "r") as f2:
    lines = set(f1.readline()) | set(f2.readline())

with open("requirements.txt" , "w") as out:
    out.writelines(sorted(lines))

