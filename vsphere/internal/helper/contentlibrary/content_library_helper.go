package contentlibrary

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/terraform-providers/terraform-provider-vsphere/vsphere/internal/helper/datastore"
	"github.com/vmware/govmomi"
	"github.com/vmware/govmomi/object"
	"github.com/vmware/govmomi/vapi/library"
	"github.com/vmware/govmomi/vapi/rest"
	"log"
)

func FromName(c *rest.Client, name string) (*library.Library, error) {
	clm := library.NewManager(c)
	ctx := context.TODO()
	lib, err := clm.GetLibraryByName(ctx, name)
	if err != nil {
		return nil, err
	}
	if lib == nil {
		return nil, fmt.Errorf("Unable to find content library (%s)", name)
	}
	return lib, err
}

func FromID(c *rest.Client, id string) (*library.Library, error) {
	clm := library.NewManager(c)
	ctx := context.TODO()
	lib, err := clm.GetLibraryByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if lib == nil {
		return nil, fmt.Errorf("Unable to find content library (%s)", id)
	}
	return lib, err
}

func CreateLibrary(c *rest.Client, name string, description string, backings []library.StorageBackings) (string, error) {
	clm := library.NewManager(c)
	ctx := context.TODO()
	lib := library.Library{
		Description: description,
		Name:        name,
		Storage:     backings,
		Type:        "LOCAL", // govmomi only supports LOCAL library creation
	}
	return clm.CreateLibrary(ctx, lib)
}

func UpdateLibrary(c *rest.Client, ol *library.Library, name string, description string, backings []library.StorageBackings) error {
	// Not currently supported in govmomi
	return nil
}

func DeleteLibrary(c *rest.Client, lib *library.Library) error {
	clm := library.NewManager(c)
	ctx := context.TODO()
	return clm.DeleteLibrary(ctx, lib)
}

func ItemFromName(c *rest.Client, l *library.Library, name string) (*library.Item, error) {
	clm := library.NewManager(c)
	ctx := context.TODO()
	fi := library.FindItem{
		LibraryID: l.ID,
		Name:      name,
	}
	items, err := clm.FindLibraryItems(ctx, fi)
	if err != nil {
		return nil, nil
	}
	if len(items) < 1 {
		return nil, fmt.Errorf("Unable to find content library item (%s)", name)
	}
	item, err := clm.GetLibraryItem(ctx, items[0])
	if err != nil {
		return nil, err
	}
	log.Print("BILLLLLLLLLLLLLLLLLLLLLLLLLLLLLLLLLL@@@@- %v", item.Type)
	return item, nil
}

func ItemFromID(c *rest.Client, id string) (*library.Item, error) {
	clm := library.NewManager(c)
	ctx := context.TODO()
	return clm.GetLibraryItem(ctx, id)
}

func CreateLibraryItem(c *rest.Client, l *library.Library, name string, desc string, source *object.VirtualMachine) (string, error) {
	clm := library.NewManager(c)
	ctx := context.TODO()
	item := library.Item{
		Description: desc,
		LibraryID:   l.ID,
		Name:        name,
		SourceID:    source.Reference().String(),
		Type:        "ovf",
	}
	return clm.CreateLibraryItem(ctx, item)
}

func UpdateLibraryItem(c *rest.Client, l *library.Library, oi *library.Item, name string, desc string) (string, error) {
	clm := library.NewManager(c)
	ctx := context.TODO()
	item := library.Item{
		Description: desc,
		LibraryID:   l.ID,
		ID:          oi.ID,
		Name:        name,
	}
	return clm.CreateLibraryItem(ctx, item)
}

func DeleteLibraryItem(c *rest.Client, item *library.Item) error {
	clm := library.NewManager(c)
	ctx := context.TODO()
	return clm.DeleteLibraryItem(ctx, item)
}

func ExpandStorageBackings(c *govmomi.Client, d *schema.ResourceData) ([]library.StorageBackings, error) {
	sb := []library.StorageBackings{}
	for _, dsId := range d.Get("storage_backing").(*schema.Set).List() {
		ds, err := datastore.FromID(c, dsId.(string))
		if err != nil {
			return nil, err
		}
		sb = append(sb, library.StorageBackings{
			DatastoreID: ds.Reference().Value,
			Type:        "DATASTORE",
		})
	}
	return sb, nil
}

func FlattenStorageBackings(sb []library.StorageBackings) []string {
	sbl := []string{}
	for _, backing := range sb {
		if backing.Type == "DATASTORE" {
			sbl = append(sbl, backing.DatastoreID)
		}
	}
	return sbl
}
