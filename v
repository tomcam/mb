# Checks a CSS file. Requires use of the 
# Nu HTML checker, which only checks HTML files.
# So you create a stub HTML file, embed the
# contents of the input file into it
# in the form of a style tag, and
# check the resulting file. So the line
# numbers refer to the HTML file, not the
# original CSS file.
# The stub HTML file, given in $OUTFILE 
# below, remains.
OUTFILE="VNU_CHECK.HTML"
# Read the file named on the command line, which is $1, 
# into the variable $INFILE.
INFILE=$(<$1)
# Create the variable $CONTENTS from this HTML, plus
# insert the contents of $1 inside a style tag. 
# Keep it all on the same line 
# so error messages are accurate
read -r -d '' CONTENTS << EOM
<!DOCTYPE html><html lang=""><head><title>$1</title><style>$INFILE</style></head><body></body></html>
EOM
## Copy all of this to the output file
echo "$CONTENTS" > $OUTFILE
echo $1:
vnu $OUTFILE


