"""
This script generates the license documentation pages.
"""

from os import system, getcwd, listdir, chmod, rename
from os.path import dirname, basename

# Directory where the scraped licenses are.
Dir = dirname(__file__) + "/../web/"

# Copy the licenses from customlicenses to licenses documentation.
Licenses = listdir(Dir + "customlicenses")
for License in Licenses:
	system("cp " + Dir + "customlicenses/" + License + " " + Dir + "docs/licenses/" + License)
Dir += "docs/licenses"

IndexPath = Dir + ".md"
IndexText = ""

# Traverse the scraped licenses and create a markdown-formatted licensing page for each one.
Licenses = listdir(Dir)
for License in Licenses:
	Path = Dir + "/" + License

	# make files writeable
	chmod(Path, 0o0600)
	# Read the license content and add it as a codeblock.
	with open(Path, "r") as File:
		TheText = "```\n" + File.read().strip() + "\n```"

	# use filename as license name, remove host prefix
	Name = License.partition("-")[2].removesuffix(".txt")
	RelName = Name
	Name += " license"

	# Generate the document, which has the following: the title, a link to the index file and the license content.
	Text = f"---\nTitle: {Name}\n---\n# {Name}"
	Text += "\n\n[Back to licenses](/docs/licenses)\n\n"
	Text += "\n" + TheText

	if RelName == "stylus-logo":
		Text += "\n\n\n#### Changes:\n- vectorized by [pabli24](https://github.com/pabli24)\n- converted to monochrome outline by [0eoc](https://userstyles.world/user/0eoc)"

	with open(Path, "w") as File:
		File.write(Text)

	# rename file
	rename(Path, Dir + "/" + RelName + ".md")

	# Add this license to the index page.
	IndexText += "\n- [" + RelName + "](/docs/licenses/" + RelName + ")\n"

# write license index. we could write to it on every license but i prefer writing big chunks
with open(IndexPath, "w") as IndexFile:
	IndexFile.write(f"---\nTitle: Third-party Licenses\n---\n# Third-party Licenses\n{IndexText}")
