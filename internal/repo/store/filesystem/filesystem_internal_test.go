package filesystem

import (
	"context"
	"errors"
	"io/fs"
	"os"
	"path"
	"testing"

	"github.com/stretchr/testify/assert"
	"gopkg.in/yaml.v2"

	"github.com/thibaultmg/clingua/internal/common"
)

var testCard = card{
	FromLanguage: "fr",
	ToLanguage:   "en",
	ID:           "myID",
}

func TestFSRepo_Get(t *testing.T) {
	assert := assert.New(t)
	t.Parallel()

	d, err := yaml.Marshal(&testCard)
	assert.Nil(err)

	tempDir, err := os.MkdirTemp("", "fsrepo_test_dir")
	assert.Nil(err)

	defer os.RemoveAll(tempDir)

	// Write valid card in dir
	err = os.WriteFile(path.Join(tempDir, "myValidCard.yaml"), d, 0o640)
	assert.Nil(err)

	// Write invalid card in dir to ensure that it is well handled by the repo
	err = os.WriteFile(path.Join(tempDir, "invalidCard.yaml"), []byte("blablabla"), 0o640)
	assert.Nil(err)

	fsrepo := New(tempDir)

	// get non existing file
	_, err = fsrepo.Get(context.Background(), "not_existing")
	assert.NotNil(err)
	assert.True(errors.Is(err, common.ErrNotFound))

	// get valid yaml file
	cardData, err := fsrepo.Get(context.Background(), testCard.ID)
	assert.Nil(err)
	assert.Equal(cardData.From.String(), testCard.FromLanguage)
}

func TestFSRepo_Delete(t *testing.T) {
	assert := assert.New(t)
	t.Parallel()

	tempDir, err := os.MkdirTemp("", "fsrepo_test_dir")
	assert.Nil(err)

	defer os.RemoveAll(tempDir)

	d, err := yaml.Marshal(&testCard)
	assert.Nil(err)

	filename := "afile.yaml"
	err = os.WriteFile(path.Join(tempDir, filename), d, 0o640)
	assert.Nil(err)

	fsrepo := New(tempDir)

	// delete non existing id
	err = fsrepo.Delete(context.Background(), "invalid")
	assert.NotNil(err)
	assert.True(errors.Is(err, common.ErrNotFound))

	// delete existing id
	err = fsrepo.Delete(context.Background(), testCard.ID)
	assert.Nil(err)

	_, err = os.ReadFile(path.Join(tempDir, filename))
	assert.NotNil(err)
	assert.ErrorIs(err, fs.ErrNotExist)
}

func TestFSRepo_Create(t *testing.T) {
	assert := assert.New(t)
	t.Parallel()

	tempDir, err := os.MkdirTemp("", "fsrepo_test_dir")
	assert.Nil(err)

	defer os.RemoveAll(tempDir)

	fsrepo := New(tempDir)

	// Create card with title
	mycard := testCard
	mycard.Title = "to ace"
	cardID, err := fsrepo.Create(context.Background(), mycard.ToEntity())
	assert.Nil(err)
	assert.True(len(cardID) > 0)

	// Create card with same id
	_, err = fsrepo.Create(context.Background(), mycard.ToEntity())
	assert.NotNil(err)
	assert.ErrorIs(err, common.ErrAlreadyExists)

	// Create card with no id
	mycard.ID = ""
	cardID, err = fsrepo.Create(context.Background(), mycard.ToEntity())
	assert.Nil(err)
	assert.Equal(cardID, mycard.ID)
}

func TestFSRepo_List(t *testing.T) {
	assert := assert.New(t)
	t.Parallel()

	tempDir, err := os.MkdirTemp("", "fsrepo_test_dir")
	assert.Nil(err)

	defer os.RemoveAll(tempDir)

	fsrepo := New(tempDir)

	// Test with no files
	cardsList, err := fsrepo.List(context.Background())
	assert.Nil(err)
	assert.Len(cardsList, 0)

	// Test with invalid card
	invalidFilePath := path.Join(tempDir, "invalid.yaml")
	err = os.WriteFile(invalidFilePath, []byte("blabla"), 0o640)
	assert.Nil(err)

	_, err = fsrepo.List(context.Background())
	assert.NotNil(err)

	err = os.Remove(invalidFilePath)
	assert.Nil(err)

	// Create files
	cardAce := testCard
	cardAce.Title = "ace"
	d, err := yaml.Marshal(cardAce)
	assert.Nil(err)
	err = os.WriteFile(path.Join(tempDir, "ace.yaml"), d, 0o640)
	assert.Nil(err)

	cardCar := testCard
	cardCar.Title = "car"
	d, err = yaml.Marshal(cardCar)
	assert.Nil(err)

	err = os.WriteFile(path.Join(tempDir, "car.yaml"), d, 0o640)
	assert.Nil(err)

	cardBoat := testCard
	cardBoat.Title = "boat"
	d, err = yaml.Marshal(cardBoat)
	assert.Nil(err)

	err = os.WriteFile(path.Join(tempDir, "boat.yaml"), d, 0o640)
	assert.Nil(err)

	// Test that files are listed
	cardsList, err = fsrepo.List(context.Background())
	assert.Nil(err)
	assert.Len(cardsList, 3)
	assert.Equal(cardAce.Title, cardsList[0].Title)
	assert.Equal(cardBoat.Title, cardsList[1].Title)
	assert.Equal(cardCar.Title, cardsList[2].Title)
}

func TestFSRepo_Cache(t *testing.T) {
	assert := assert.New(t)
	t.Parallel()

	tempDir, err := os.MkdirTemp("", "fsrepo_test_dir")
	assert.Nil(err)

	defer os.RemoveAll(tempDir)

	fsrepo := New(tempDir)

	_, err = fsrepo.Create(context.Background(), testCard.ToEntity())
	assert.Nil(err)

	// Test that cache is filled with new entry
	_, err = fsrepo.Get(context.Background(), testCard.ID)
	assert.Nil(err)

	// Delete entry and check that it is removed from cache
	err = fsrepo.Delete(context.Background(), testCard.ID)
	assert.Nil(err)

	_, err = fsrepo.Get(context.Background(), testCard.ID)
	assert.NotNil(err)
}
