package filesystem

import (
	"context"
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path"
	"strings"

	"gopkg.in/yaml.v2"

	"github.com/thibaultmg/clingua/internal/common"
	"github.com/thibaultmg/clingua/internal/entity"
)

const (
	yamlExtension    = ".yaml"
	filesPermissions = 0o640
)

type FSRepo struct {
	root string
}

func New(root string) *FSRepo {
	if !path.IsAbs(root) {
		panic("invalid root path")
	}

	return &FSRepo{
		root: root,
	}
}

// Get returns the Card having the ID id, which is its file path.
func (f *FSRepo) Get(ctx context.Context, id string) (entity.Card, error) {
	fileData, err := os.ReadFile(path.Join(f.root, id+yamlExtension))
	// fileData, err := fs.ReadFile(f.fs, id)
	if err != nil {
		if errors.Is(err, fs.ErrNotExist) {
			return entity.Card{}, common.ErrNotFound
		}

		return entity.Card{}, common.NewErrInternalError(err)
	}

	var cardData card

	err = yaml.Unmarshal(fileData, &cardData)
	if err != nil {
		return entity.Card{}, common.NewErrInternalError(err)
	}

	return cardData.ToEntity(), nil
}

func (f *FSRepo) Create(ctx context.Context, ecard entity.Card) (string, error) {
	card := entityToCard(&ecard)

	cardData, err := yaml.Marshal(&card)
	if err != nil {
		return "", err
	}

	fileName := card.Title
	fileName = strings.TrimSpace(fileName)
	fileName = strings.ToLower(fileName)
	fileName = strings.Join(strings.Fields(fileName), "_")

	if len(fileName) == 0 {
		fileName = "no_title"
	}

	var (
		cardFile               *os.File
		antiCollisionExtension string
		counter                int
	)

main:
	for {
		cardFile, err = os.OpenFile(
			path.Join(f.root, fileName+antiCollisionExtension+yamlExtension), os.O_WRONLY|os.O_CREATE|os.O_EXCL, filesPermissions)

		switch {
		case errors.Is(err, fs.ErrExist):
			counter++
			if counter > 0 {
				antiCollisionExtension = fmt.Sprintf("(%d)", counter)
			}
		case err != nil:
			return fileName, err
		default:
			defer cardFile.Close()
			fileName += antiCollisionExtension

			break main
		}
	}

	_, err = cardFile.Write(cardData)
	if err != nil {
		return fileName, err
	}

	return fileName, nil
}

func (f *FSRepo) Delete(ctx context.Context, id string) error {
	filePath := path.Join(f.root, id+yamlExtension)

	err := os.Remove(filePath)
	if err != nil {
		if errors.Is(err, fs.ErrNotExist) {
			return common.ErrNotFound
		}

		return common.NewErrInternalError(err)
	}

	return nil
}

func (f *FSRepo) List(ctx context.Context) ([]entity.Card, error) {
	dirEntries, err := os.ReadDir(f.root)
	if err != nil {
		return []entity.Card{}, err
	}

	ret := make([]entity.Card, 0, len(dirEntries))

	for _, e := range dirEntries {
		if e.IsDir() {
			continue
		}

		if !strings.HasSuffix(e.Name(), yamlExtension) {
			continue
		}

		fileData, err := os.ReadFile(path.Join(f.root, e.Name()))
		if err != nil {
			return []entity.Card{}, err
		}

		var repoCard card

		err = yaml.Unmarshal(fileData, &repoCard)
		if err != nil {
			return []entity.Card{}, err
		}

		ret = append(ret, repoCard.ToEntity())
	}

	return ret, nil
}

func (f *FSRepo) Search(name string) error {
	return nil
}
