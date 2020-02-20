package vsphere

import (
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/terraform-providers/terraform-provider-vsphere/vsphere/internal/helper/contentlibrary"
	"github.com/terraform-providers/terraform-provider-vsphere/vsphere/internal/helper/virtualmachine"
)

func resourceVSphereContentLibraryItem() *schema.Resource {
	return &schema.Resource{
		Create: resourceVSphereContentLibraryItemCreate,
		Delete: resourceVSphereContentLibraryItemDelete,
		Read:   resourceVSphereContentLibraryItemRead,
		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The name of the content library item.",
			},
			"library_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "ID of the content library to contain item",
			},
			"description": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: "Optional description of the content library item.",
			},
			"source_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "ID of source VM of content library item.",
			},
		},
	}
}

func resourceVSphereContentLibraryItemCreate(d *schema.ResourceData, meta interface{}) error {
	vc := meta.(*VSphereClient).vimClient
	rc := meta.(*VSphereClient).restClient

	source, err := virtualmachine.FromUUID(vc, d.Get("source_id").(string))
	if err != nil {
		return err
	}

	lib, err := contentlibrary.FromID(rc, d.Get("library_id").(string))
	if err != nil {
		return err
	}

	id, err := contentlibrary.CreateLibraryItem(rc, lib, d.Get("name").(string), d.Get("description").(string), source)
	if err != nil {
		return err
	}
	d.SetId(id)
	return resourceVSphereContentLibraryItemRead(d, meta)
}

func resourceVSphereContentLibraryItemDelete(d *schema.ResourceData, meta interface{}) error {
	rc := meta.(*VSphereClient).restClient
	item, err := contentlibrary.ItemFromID(rc, d.Id())
	if err != nil {
		return err
	}
	return contentlibrary.DeleteLibraryItem(rc, item)
}

func resourceVSphereContentLibraryItemRead(d *schema.ResourceData, meta interface{}) error {
	rc := meta.(*VSphereClient).restClient
	item, err := contentlibrary.ItemFromID(rc, d.Id())
	if err != nil {
		return err
	}
	d.Set("name", item.Name)
	d.Set("description", item.Description)
	return nil
}
