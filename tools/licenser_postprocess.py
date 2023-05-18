"""
This script generates the license documentation pages.
"""

from os import system, getcwd, listdir, chmod, rename
from os.path import join, dirname, basename
from shutil import copyfile

# Optionally specify the changes for GPL-licensed dependencies.
ChangedDependencies = {
	"stylus-logo": "\n- vectorized by [pabli24](https://github.com/pabli24)\n- converted to monochrome outline by [0eoc](https://userstyles.world/user/0eoc)",
}

# Directory where the missing licenses are.
CustomDir = join(dirname(__file__), "../web/customlicenses" )

# Directory where the scraped licenses are.
Dir = join(dirname(__file__), "../web/docs/licenses")

# Copy the licenses from customlicenses to licenses documentation.
Licenses = listdir(CustomDir)
for License in Licenses:
	copyfile(join(CustomDir, License), join(Dir, License))

IndexPath = Dir + ".md"
Index = []

# Traverse the scraped licenses and create a markdown-formatted licensing page for each one.
Licenses = listdir(Dir)
for License in Licenses:
	Path = Dir + "/" + License

	# Make files writeable, as they may not be.
	chmod(Path, 0o0600)

	# Read the license content and add it as a codeblock.
	with open(Path, "r") as File:
		LicenseText = File.read().strip()

	# Use filename as license name, remove host prefix.
	Name = License.partition("-")[2].removesuffix(".txt")
	RelName = Name
	Name += " license"

	# Generate the document, which has the following: the title, a link to the index file and the license content.
	Text = f"---\nTitle: {Name}\n---\n# {Name}"
	Text += "\n\n[Back to licenses](/docs/licenses)\n\n"
	Text += "\n```\n" + LicenseText + "\n```"

	if RelName in ChangedDependencies:
		Text += "\n\n\n#### Changes:" + ChangedDependencies[RelName]

	with open(Path, "w") as File:
		File.write(Text)

	# Change the file extension to the correct one.
	rename(Path, Dir + "/" + RelName + ".md")

	# Add this license to the index.
	Index.append(RelName)

# Generate the license index.
Index.sort()
IndexText = "---\nTitle: Third-party Licenses\n---\n# Third-party Licenses\n"
for License in Index:
	IndexText += f"\n- [{License}](/docs/licenses/{License})\n"

# Write the index. We could write to it on every license, but for performance reasons, it's preferred to write to it all at once.
with open(IndexPath, "w") as IndexFile:
	IndexFile.write(IndexText)
