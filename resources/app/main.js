jQuery(document).ready(() => {

    let downloadLink = $("#download-link")
    let downloadLinkAlert = $("#download-link-alert")
    let downloadBtn = $("#download-btn")

    let video = $("#video")
    let videoTitle = $("#video-title")
    let videoBvid = $("#video-bvid")
    let videoDescription = $("#video-description")
    let videoPic = $("#video-pic")

    downloadLinkAlert.hide()
    video.hide()

    // on user changes the value of the download link
    downloadLink.on("input propertychange", event => {
        event.preventDefault()
        let download_link = event.target.value
        if (download_link.length === 0) {
            downloadLinkAlert.hide()
            video.hide()
            return
        }

        astilectron.sendMessage({"id": "downloader.url_update", "payload": download_link}, message => {
            // after the download link changes is handled
            if (message.id === "downloader.error") {
                // the download link's format is invalid
                video.hide()
                downloadLinkAlert.text(message.payload)
                downloadLinkAlert.show()
            } else if (message.id === "downloader.url_update") {
                // the download link's format is valid, display the video info
                downloadLinkAlert.hide()
                videoTitle.text(message.payload.title)
                videoBvid.text(message.payload.bvid)
                videoDescription.text(message.payload.desc)
                videoPic.attr("href", message.payload.pic)
                video.show()
            }
        })
    })

    // on user submit a download request by clicking the button
    downloadBtn.click(() => {
        astilectron.sendMessage({"id": "downloader.download_request", "payload": downloadLink.val()}, message => {
            if (message.id === "downloader.download_request")
                alert("下载成功，请检查本地用户home目录！")
            else
                alert("下载失败：" + JSON.stringify(message.payload))
        })
    })

    // This will wait for the astilectron namespace to be ready
    document.addEventListener('astilectron-ready', function () {
        // This will listen to messages sent by GO
        astilectron.onMessage((message) => {
            if (message.id === "downloader.download_update") {
                const {cur_length, total_length} = message.payload
                const percentage = cur_length / total_length * 100
                $(".progress-bar").css("width", `${percentage}%`).text(`${cur_length} / ${total_length}`);
            }
        });
    })
})


