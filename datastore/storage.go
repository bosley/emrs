package datastore

import (
	"errors"
	c "github.com/ostafen/clover/v2"
	"github.com/ostafen/clover/v2/query"
	"log/slog"
)

type storage struct {
	users  *c.DB
	assets *c.DB
}

func (ds *storage) Close() {
	ds.users.Close()
	ds.assets.Close()
}

func (ds *storage) AddAsset(asset Asset) bool {

	results, err := ds.assets.FindAll(query.NewQuery(assetDb).Where(query.Field("id").Eq(asset.Id)))
	if err != nil {
		slog.Error("error arose while attempting to find asset to judge uniqueness", "id", asset.Id, "error", err.Error())
		return false
	}

	if len(results) != 0 {
		slog.Error("asset with given id already exists", "id", asset.Id)
		return false
	}

	did, err := ds.assets.InsertOne(assetDb, asset.ToDoc())
	if err != nil {
		slog.Error("failed to insert asset into assets collection", "id", asset.Id, "error", err.Error())
		return false
	}

	slog.Debug("added asset to database", "id", asset.Id, "doc-id", did)

	return true
}

func (ds *storage) RemoveAsset(id string) bool {
	err := ds.assets.Delete(query.NewQuery(assetDb).Where(query.Field("id").Eq(id)))
	if err != nil {
		if errors.Is(err, c.ErrIndexNotExist) {
			slog.Info("asset for removal was not found in the database", "id", id)
		} else {
			slog.Error("error attempting to remove asset", "id", id, "error", err.Error())
			return false
		}
	}
	return true
}

func (ds *storage) UpdateAsset(asset Asset) bool {

	results, err := ds.assets.FindAll(query.NewQuery(assetDb).Where(query.Field("id").Eq(asset.Id)))
	if err != nil {
		slog.Error("error arose while attempting to find asset to judge uniqueness", "id", asset.Id, "error", err.Error())
		return false
	}

	if len(results) == 0 {
		slog.Error("asset with given id does not exist", "id", asset.Id)
		return false
	}

	x := make(map[string]interface{})
	x["displayname"] = asset.DisplayName

	err = ds.assets.Update(query.NewQuery(assetDb).Where(query.Field("id").Eq(asset.Id)), x)
	if err != nil {
		slog.Error("error updating asset", "id", asset.Id, "error", err.Error())
		return false
	}
	return true
}

func (ds *storage) GetAssets() []Asset {

	result := make([]Asset, 0)

	docs, err := ds.assets.FindAll(query.NewQuery(assetDb))
	if err != nil {
		slog.Error("failed to retrieve assets", "error", err.Error())
		return result
	}

	for _, doc := range docs {
		t := &Asset{}
		doc.Unmarshal(t)
		result = append(result, *t)
	}
	return result
}

func (ds *storage) UpdateOwnerUiKey(key string) bool {

	x := make(map[string]interface{})
	x["uikey"] = key

	err := ds.assets.Update(query.NewQuery(usersDb).Where(query.Field("ring").Eq(RingOne)), x)
	if err != nil {
		slog.Error("error updating ui key", "error", err.Error())
		return false
	}
	return true
}

func (ds *storage) GetOwner() (User, error) {

	docs, err := ds.users.FindAll(query.NewQuery(usersDb).Where(query.Field("ring").Eq(RingOne)))
	if err != nil {
		slog.Error("failed to retrieve ring-1 user", "error", err.Error())
		return User{}, err
	}

	if len(docs) != 1 {
		slog.Warn("server misconfiguration detected")
		slog.Error("failed to retrieve owner, multiple were found", "n", len(docs))
		return User{}, nil
	}

	t := &User{}
	docs[0].Unmarshal(t)
	return *t, nil
}

func (ds *storage) UpdateOwner(owner User) bool {

	docs, err := ds.users.FindAll(query.NewQuery(usersDb).Where(query.Field("ring").Eq(RingOne)))
	if err != nil {
		slog.Error("failed to retrieve ring-1 user", "error", err.Error())
		return false
	}

	if len(docs) != 1 {
		slog.Warn("server misconfiguration detected")
		slog.Error("failed to retrieve owner, multiple were found", "n", len(docs))
		return false
	}

	owner.Ring = RingOne

	err = ds.users.Update(query.NewQuery(usersDb).Where(query.Field("ring").Eq(RingOne)), owner.ToDoc().AsMap())
	if err != nil {
		slog.Error("Error updating owner information", "error", err.Error())
		return false
	}

	return true
}
