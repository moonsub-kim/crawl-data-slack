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
			});
		}
		parse_dto();
		`,
	},
}
