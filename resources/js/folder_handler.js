const videoList = document.querySelector(".video-list");
const player = document.getElementById("videoPlayer");

window.currentVideos = [];
window.currentIndex = -1;

document.getElementById("folderInput").addEventListener("change", function (e) {
    videoList.innerHTML = "";

    const files = Array.from(e.target.files)
        .filter(file => file.type.startsWith("video/"));

    if (files.length === 0) {
        alert("No video files found!");
        return;
    }

    window.currentVideos = files;

    files.forEach((file, index) => {
        const li = document.createElement("li");
        li.textContent = file.webkitRelativePath || file.name;

        li.onclick = () => {
            window.currentIndex = index;
            const url = URL.createObjectURL(file);
            player.src = url;

            document.querySelectorAll(".video-list li").forEach(el => el.classList.remove("active"));
            li.classList.add("active");

            // if (window.sendControl) {
            //     window.sendControl("jump", { index: index });
            // }
        };

        videoList.appendChild(li);
    });
});

function toggleSidebar() {
    document.getElementById("sidebar").classList.toggle("collapsed");
}

document.getElementById("folderInput").addEventListener("click", () => {
    const url = URL.createObjectURL(file);
    player.src = url;
});