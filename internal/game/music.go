package game

import (
  rl "github.com/gen2brain/raylib-go/raylib"
)

var music rl.Music

func InitMusic() {
  rl.InitAudioDevice()
  music = rl.LoadMusicStream("assets/sounds/Paulsfliegekurs.wav")
  rl.PlayMusicStream(music)
}

func PlayMusic() {

  rl.UpdateMusicStream(music)
}

func QuitMusic() {
  rl.UnloadMusicStream(music)
  rl.CloseAudioDevice()
}
