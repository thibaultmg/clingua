package filesystem

import (
	"context"
	"errors"
	"io/fs"
	"os"
	"path"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/thibaultmg/clingua/internal/common"
	"gopkg.in/yaml.v2"
)

var testCard = card{
	FromLanguage: "fr",
}

func TestFSRepo_Get(t *testing.T) {
	assert := assert.New(t)

	d, err := yaml.Marshal(&testCard)
	assert.Nil(err)

	tempDir, err := os.MkdirTemp("", "fsrepo_test_dir")
	assert.Nil(err)

	defer os.RemoveAll(tempDir)

	// Write valid card in dir
	fileID := "myValidCard"
	err = os.WriteFile(path.Join(tempDir, fileID+".yaml"), d, 0664)
	assert.Nil(err)

	// Write invalid card in dir
	invalidFileID := "invalidCard"
	err = os.WriteFile(path.Join(tempDir, invalidFileID+".yaml"), []byte("blablabla"), 0664)
	assert.Nil(err)

	fsrepo := New(tempDir)

	// get non existing file
	_, err = fsrepo.Get(context.Background(), "not_existing")
	assert.NotNil(err)
	assert.True(errors.Is(err, common.ErrNotFound))

	// get invalid yaml file
	_, err = fsrepo.Get(context.Background(), invalidFileID)
	assert.NotNil(err)
	assert.True(errors.As(err, &common.ErrInternalError{}))

	// get valid yaml file
	cardData, err := fsrepo.Get(context.Background(), fileID)
	assert.Nil(err)
	assert.Equal(cardData.From.String(), testCard.FromLanguage)
}

func TestFSRepo_Delete(t *testing.T) {
	assert := assert.New(t)

	tempDir, err := os.MkdirTemp("", "fsrepo_test_dir")
	assert.Nil(err)

	defer os.RemoveAll(tempDir)

	fileID := "myFile"
	err = os.WriteFile(path.Join(tempDir, fileID+".yaml"), []byte("blablablagarbage"), 0664)
	assert.Nil(err)

	fsrepo := New(tempDir)

	// delete non existing file
	err = fsrepo.Delete(context.Background(), "invalid")
	assert.NotNil(err)
	assert.True(errors.Is(err, common.ErrNotFound))

	// delete existing file
	err = fsrepo.Delete(context.Background(), fileID)
	assert.Nil(err)
	_, err = os.ReadFile(path.Join(tempDir, fileID+".yaml"))
	assert.NotNil(err)
	assert.ErrorIs(err, fs.ErrNotExist)
}

func TestFSRepo_Create(t *testing.T) {
	assert := assert.New(t)

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

	// Create card with same title
	cardID, err = fsrepo.Create(context.Background(), mycard.ToEntity())
	assert.Nil(err)
	assert.True(len(cardID) > 0)

	// Create card with no title
	mycard.Title = ""
	cardID, err = fsrepo.Create(context.Background(), mycard.ToEntity())
	assert.Nil(err)
	assert.True(len(cardID) > 0)
}

func TestFSRepo_List(t *testing.T) {
	assert := assert.New(t)

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
	os.WriteFile(invalidFilePath, []byte("blabla"), 0664)
	cardsList, err = fsrepo.List(context.Background())
	assert.NotNil(err)

	err = os.Remove(invalidFilePath)
	assert.Nil(err)

	// Create files
	cardAce := testCard
	cardAce.Title = "ace"
	d, err := yaml.Marshal(cardAce)
	assert.Nil(err)
	os.WriteFile(path.Join(tempDir, "ace.yaml"), d, 0664)

	cardCar := testCard
	cardCar.Title = "car"
	d, err = yaml.Marshal(cardCar)
	assert.Nil(err)
	os.WriteFile(path.Join(tempDir, "car.yaml"), d, 0664)

	cardBoat := testCard
	cardBoat.Title = "boat"
	d, err = yaml.Marshal(cardBoat)
	assert.Nil(err)
	os.WriteFile(path.Join(tempDir, "boat.yaml"), d, 0664)

	// Test that files are listed
	cardsList, err = fsrepo.List(context.Background())
	assert.Nil(err)
	assert.Len(cardsList, 3)
	assert.Equal(cardAce.Title, cardsList[0].Title)
	assert.Equal(cardBoat.Title, cardsList[1].Title)
	assert.Equal(cardCar.Title, cardsList[2].Title)
}
