
all: bdownloader

bdownloader:
	rm -rf output
	mkdir output
	go build -o output/bdownloader github.com/hanfa/bili-downloader
	cp -r resources output

clean:
	rm -rf output vendor


