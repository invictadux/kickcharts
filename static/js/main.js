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

function SelectSortAction(s) {
	const url = new URL(window.location.href);
	url.searchParams.delete('page');
	url.searchParams.set('sort', s.dataset.value);
	window.location.href = url.toString();
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

function GenerateBarChartOptions(dates, values, title, color) {
	var option = {
		backgroundColor: '#1c1c22',
		title: [
			{
				text: title,
				left: 'left',
				top: 'top',
				textStyle: {
					fontFamily: 'Karla',
					fontSize: 12,
					fontWeight: '400',
					color: '#dedede',
				}
			},
			{
				text: 0,
				right: 'right',
				top: 'top',
				textStyle: {
					fontFamily: 'Karla',
					fontSize: 13,
					fontWeight: '400',
					color: '#dedede',
				}
			},
		],
		tooltip: {
			trigger: 'axis',
			backgroundColor: 'ffffff00',
			borderWidth: 0,
			textStyle: {
				color: '#ffffff',
				fontSize: 14
			},
			axisPointer: {
				type: 'shadow'
			},
			formatter: function(params) {
				option.title[1].text = params[0].value.toLocaleString();
				return `${params[0].name}: ${params[0].value.toLocaleString()}`;
			},
			position: function (point, params, dom, rect, size) {
				return [size.viewSize[0] - size.contentSize[0], 0];
			}
		},
		grid: {
			left: '1%',
			right: '1%',
			top: '20%',
			bottom: '0%',
			containLabel: true,
		},
		xAxis: {
			type: 'category',
			data: dates,
			axisLabel: { show: false },
			axisLine: { show: false },
			axisTick: { show: false }, 
		},
		yAxis: {
			type: 'value',
			splitLine: {
				show: false,
			},
			axisLabel: { show: false },
			axisLine: { show: false },
			axisTick: { show: false }, 
		},
		series: [
			{
				name: 'Direct',
				data: values,
				type: 'bar',
				barWidth: '90%',
				itemStyle: {
                    color: color,
                }
			}
		]
	};

	return option;
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