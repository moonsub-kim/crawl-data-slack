package eomisae

type target struct {
	name   string
	url    string
	script string
}

var targets map[string]target = map[string]target{
	"raffle": {
		name: "raffle",
		url:  "https://eomisae.co.kr/dr",
		script: `
		function parse_dto() {
			return JSON.stringify({
				'post': document.URL,
				'name': document.querySelector('h2 > a.pjax').innerText,
				'url': document.querySelector('td.extra_url > a').href,
				'content': '',
			});
		}
		parse_dto();
		`,
	},
	"fashion": { // document.querySelector('div.card_content > h3 > a.pjax').innerHTML.includes('7 레벨 미만') === true 면 continue
		name: "fashion",
		url:  "https://eomisae.co.kr/os",
		script: `
		function parse_dto() {
			return JSON.stringify({
				'post': document.URL,
				'name': document.querySelector('h2 > a.pjax').innerText,
				'url': document.querySelector('td.extra_url > a').href,
				'content': document.querySelector('div.xe_content').innerText,
			});
		}
		parse_dto();
		`,
	},
}
