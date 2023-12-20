# Gomagnon

![Agents des deux types se déplaçant sur la carte, resources diverses, et illustration de la génération du monde.](assets%2Fimages%2Fglobal.png "Screenshot de l'interface du projet.")


Simulation multi-agent écrite en Go d'un ecosystème préhistorique. Créée par Adrien Simon, Quentin Fitte-Rey, Jean Lescieux, et Raphael Quintaneiro, dans le cadre du cours IA04 - Système Multi-Agents de l'UTC, Compiègne, France.

Ce projet est une tentative d'évaluation des raisons qui ont mené Néanderthal à perdre le jeu de l'occupation de la terre face à Sapiens.

---

## Installation

Clonez le projet en local et lancez le grâce à la commande `go run`. 
```shell
git clone https://github.com/adrsimon/gomagnon
cd gomagnon
go run .
```

### Configuration

La simulation est gérée par un fichier de configuration situé dans le dossier `/settings`. 
```json
{
  "agents": {
    "initialNumber": 30
  },
  "world": {
    "seed": 2021,
    "type": "continent",
    "resources": {
      "maxAnimals": 40,
      "maxFruits": 40,
      "maxWoods": 30,
      "maxRocks": 30
    },
    "size": {
      "x": 46,
      "y": 41
    }
  }
}
```
`initialNumber` : nombre d'agents à l'origine du monde<br>
`seed` : la graîne de génération du monde, qui vous permettra de pouvoir reproduire la même solution plusieurs fois d'affilée.<br>
`type` : type du monde, doit être un de `ìsland|continent`, island générera une grande île entourée d'eau, continent génerera un continent avec plusieurs lacs.<br>
`maxXXX` : quantité maximale de la ressource XXX simulaténement disponible sur la carte<br>
`size` : taille du monde, X correspond à la largeur, Y à la hauteur.