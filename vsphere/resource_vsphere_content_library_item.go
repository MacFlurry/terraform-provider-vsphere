package vsphere

import (
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/terraform-providers/terraform-provider-vsphere/vsphere/internal/helper/contentlibrary"
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
			"file_url": {
				Type:        schema.TypeSet,
				Required:    true,
				ForceNew:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "ID of source VM of content library item.",
			},
			"type": {
				Type: schema.TypeString,
				Default: "ovf",
				Optional: true,
				Description: "Type of content library item.",
			},
		},
	}
}

func resourceVSphereContentLibraryItemCreate(d *schema.ResourceData, meta interface{}) error {
	rc := meta.(*VSphereClient).restClient

	lib, err := contentlibrary.FromID(rc, d.Get("library_id").(string))
	if err != nil {
		return err
	}

	files := d.Get("file_url").(*schema.Set)

	id, err := contentlibrary.CreateLibraryItem(rc, lib, d.Get("name").(string), d.Get("description").(string), d.Get("type").(string), files.List())
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
