package player

import "os"
import "strings"
import "errors"
import "math/rand"
import "sort"

var NoMoreSongs error = errors.New("No more songs")

func shuffle(xs []Song, i int) int {
	x := xs[i]
	for j := 0; j < len(xs); j++ {
		r := rand.Intn(len(xs))
		xs[j], xs[r] = xs[r], xs[j]
		if xs[j] == x {
			i = j
		}
		if xs[r] == x {
			i = r
		}
	}
	return i
}

func remove(i int, xs []Song) []Song {
	return append(xs[:i], xs[i+1:]...)
}

func mod(r int, m int) int {
	t := r % m
	if t < 0 {
		t += m
	}
	return t
}

func FindSongs(dir string) (songs []Song, err error) {
	f, err := os.Open(dir)
	if err != nil {
		return
	}
	files, err := f.Readdirnames(-1)
	if err != nil {
		return
	}
	songs = []Song{}
	for _, name := range files {
		if strings.HasSuffix(name, ".mp3") {
			songs = append(songs, Song(name))
		}
	}
	return
}

type Player struct {
	Shuffle bool
	Repeat  bool
	Speaker *Speaker
	Songs   []Song
	curr    int
}

func NewPlayer(songs []Song) *Player {
	p := &Player{
		Shuffle: false,
		Repeat:  false,
		Speaker: NewSpeaker(),
		Songs:   songs,
	}
	p.sort()
	return p
}

func (p *Player) ToggleRepeat() {
	p.Repeat = !p.Repeat
}

func (p *Player) sort() {
	sort.Slice(p.Songs, func(i, j int) bool {
		return string(p.Songs[i]) < string(p.Songs[j])
	})
}

func (p *Player) ToggleShuffle() {
	p.Shuffle = !p.Shuffle
	if p.Shuffle {
		p.curr = shuffle(p.Songs, p.curr)
	} else {
		song, err := p.Song()
		p.sort()
		if err == nil {
			p.SetCurrent(song)
		}
	}
}

func (p *Player) Song() (Song, error) {
	return p.Peek(0)
}

func (p *Player) Peek(i int) (Song, error) {
	j := mod(p.curr+i, len(p.Songs))
	if len(p.Songs) == 0 {
		return Song(""), NoMoreSongs
	}
	return p.Songs[j], nil
}

func (p *Player) Next(i int, force bool) (Song, error) {
	if !p.Repeat || force {
		p.curr = mod(p.curr+i, len(p.Songs))
	}
	return p.Song()
}

func (p *Player) Remove() {
	p.Songs = remove(p.curr, p.Songs)
}

func (p *Player) Toggle() {
	p.Speaker.Toggle()
}

func (p *Player) SetCurrent(s Song) {
	for i := 0; i < len(p.Songs); i++ {
		if p.Songs[i] == s {
			p.curr = i
			break
		}
	}
}
