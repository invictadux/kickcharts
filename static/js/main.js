
function GetChannels(page) {
    page.canLoad = false;

	var http = new XMLHttpRequest();
	var url = new URL(window.location.origin + "/api/v1/channels");

	url.searchParams.set("page", page['page']);

	http.onreadystatechange = function() {
	    if (this.readyState == 4 && this.status == 200) {
	        var data = JSON.parse(this.responseText);
			RenderChannels("channels", data);
            
            if (data.length == 30) {
                page.page += 1;
                page.canLoad = true;
            }
	    }
	};

	http.open("GET", url, true);
	http.send();
}

function GetClips(page) {
    page.canLoad = false;

	var http = new XMLHttpRequest();
	var url = new URL(window.location.origin + "/api/v1/clips");

	url.searchParams.set("page", page['page']);

	http.onreadystatechange = function() {
	    if (this.readyState == 4 && this.status == 200) {
	        var data = JSON.parse(this.responseText);
			RenderClips("clips", data);
            
            if (data.length == 30) {
                page.page += 1;
                page.canLoad = true;
            }
	    }
	};

	http.open("GET", url, true);
	http.send();
}

async function FetchChartStats(s, t) {
    const url = new URL(window.location.origin + "/api/v1/chart/stats");
    url.searchParams.set("s", s);
    url.searchParams.set("t", t);

    try {
        const response = await fetch(url);
        
        if (response.ok) {
            const data = await response.json();
			return data;
        } else {
            console.error('Error fetching chart stats:', response.status, response.statusText);
			return null;
        }
    } catch (error) {
        console.error('Error occurred while fetching chart stats:', error);
		return null;
    }
}

function PlayClip(src) {
	var elem = document.createElement('div');
	elem.setAttribute('class', 'video-player');
	elem.innerHTML = `<div class="player-card">
		<svg xmlns="http://www.w3.org/2000/svg" onclick="CloseClip()" height="20" width="20" viewBox="0 0 384 512"><path d="M342.6 150.6c12.5-12.5 12.5-32.8 0-45.3s-32.8-12.5-45.3 0L192 210.7 86.6 105.4c-12.5-12.5-32.8-12.5-45.3 0s-12.5 32.8 0 45.3L146.7 256 41.4 361.4c-12.5 12.5-12.5 32.8 0 45.3s32.8 12.5 45.3 0L192 301.3 297.4 406.6c12.5 12.5 32.8 12.5 45.3 0s12.5-32.8 0-45.3L237.3 256 342.6 150.6z"/></svg>
		<video id='video-player' width="640" height="360" class="video-js vjs-theme-forest vjs-16-9 vjs-fluid" controls>
			<source src="${src}" type="application/x-mpegURL">
		</video>
	</div>`;

	document.body.appendChild(elem);

	var player = videojs('video-player');

	player.src({
	   src: src,
	   type: 'application/x-mpegURL'
	});

	player.play();
}

function CloseClip() {
	var oldplayer = document.getElementsByClassName('video-player')[0];
	videojs(document.getElementById('video-player')).dispose();
	oldplayer.remove();
}

function GetCategories(page) {
    page.canLoad = false;

	var http = new XMLHttpRequest();
	var url = new URL(window.location.origin + "/api/v1/categories");

	url.searchParams.set("page", page['page']);

	http.onreadystatechange = function() {
	    if (this.readyState == 4 && this.status == 200) {
	        var data = JSON.parse(this.responseText);
			RenderCategories("categories", data);
            
            if (data.length == 30) {
                page.page += 1;
                page.canLoad = true;
            }
	    }
	};

	http.open("GET", url, true);
	http.send();
}

function RenderChannels(target, data) {
	var container = document.getElementById(target);

	for (const channel of data) {
		container.innerHTML += `
        <div class="channel">
            <a href="/channel/${channel['slug']}"><img src="${channel['picture']}"></a>
            <a href="/channel/${channel['slug']}"><h3>${channel['username']}</h3></a>
            <p>${channel['peak_viewers']} peak viewers</p>
        </div>`;
	}
}

function RenderCategories(target, data) {
	var container = document.getElementById(target);

	for (const category of data) {
		container.innerHTML += `
        <div class="category">
            <a href="/category/${category['slug']}"><img src="${category['banner']}"></a>
            <a href="/category/${category['slug']}"><h3>${category['name']}</h3></a>
            <p>${category['peak_viewers']} peak viewers</p>
        </div>`;
	}
}

function RenderClips(target, data) {
	var container = document.getElementById(target);

	for (const channel of data) {
		container.innerHTML += `
		<clip>
			<div class="img-container" onclick="PlayClip('${channel['url']}')">
				<img src="${channel['thumbnail']}" alt="${channel['title']}" loading="lazy">
				<div class="duration">${channel['duration']}s</div>
				<div class="view-count">${channel['views']} views</div>
				<div class="created-date">${channel['created_at']}</div>
			</div>
			<div class="clip-content">
			<data>
				<h3>${channel['title']}</h3>
				<a href="/channel/${channel['channel']}">${channel['channel']}</a><a href="/category/${channel['category']}">${channel['category']}</a>
			</data>
			</div>
		</clip>`;
	}
}

for (const dropdown of document.querySelectorAll(".custom-select-wrapper")) {
    dropdown.addEventListener('click', function () {
        this.querySelector('.custom-select').classList.toggle('open');
    })
}

for (const option of document.querySelectorAll(".custom-option")) {
    if (option.classList.contains('selected')) {
        option.closest('.custom-select').querySelector('.custom-select__trigger span').textContent = option.textContent;
    }

    option.addEventListener('click', function () {
        this.closest('.custom-select').querySelector('.custom-select__trigger span').textContent = this.textContent;
            
        var wrapper = this.closest('.custom-select-wrapper');

        for (const option of this.closest('.custom-select').querySelectorAll(".custom-option")) {
            option.classList.remove('selected');
        }

        this.classList.add('selected');
        wrapper.dataset.value = this.dataset.value;
        wrapper.dispatchEvent(new CustomEvent("selectChange", {
            detail: {
                value: this.textContent,
                item: option,
            },
        }));
    })
}

window.addEventListener('click', function (e) {
    for (const select of document.querySelectorAll('.custom-select')) {
        if (!select.contains(e.target)) {
            select.classList.remove('open');
        }
    }
});