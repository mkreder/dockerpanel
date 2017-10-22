> install_deps.sh
for url in $(go list -f '{{ join .Deps "\n" }}' | egrep "^golang|^git" ); do
	echo "go get $url " >> install_deps.sh
done
