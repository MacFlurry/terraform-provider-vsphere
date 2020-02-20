package vsphere

import (
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/terraform-providers/terraform-provider-vsphere/vsphere/internal/helper/contentlibrary"
	"log"
)

func resourceVSphereContentLibrary() *schema.Resource {
	return &schema.Resource{
		Create: resourceVSphereContentLibraryCreate,
		Delete: resourceVSphereContentLibraryDelete,
		Read:   resourceVSphereContentLibraryRead,
		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The name of the content library.",
			},
			"description": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: "Optional description of the content library.",
			},
			"storage_backing": {
				Type:        schema.TypeSet,
				Required:    true,
				ForceNew:    true,
				Description: "The name of the content library.",
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
		},
	}
}

func resourceVSphereContentLibraryRead(d *schema.ResourceData, meta interface{}) error {
	c := meta.(*VSphereClient).restClient
	lib, err := contentlibrary.FromName(c, d.Get("name").(string))
	if err != nil {
		return err
	}
	d.SetId(lib.ID)
	sb := contentlibrary.FlattenStorageBackings(lib.Storage)
	d.Set("name", lib.Name)
	d.Set("description", lib.Description)
	err = d.Set("storage_backing", sb)
	if err != nil {
		return err
	}
	return nil
}

func resourceVSphereContentLibraryCreate(d *schema.ResourceData, meta interface{}) error {
	vc := meta.(*VSphereClient).vimClient
	rc := meta.(*VSphereClient).restClient
	backings, err := contentlibrary.ExpandStorageBackings(vc, d)
	log.Print("BILLLLLLLLLLLLLLLLLLLLLLLLLLL - %v", backings)
	if err != nil {
		return err
	}
	id, err := contentlibrary.CreateLibrary(rc, d.Get("name").(string), d.Get("description").(string), backings)
	if err != nil {
		return err
	}
	d.SetId(id)
	return resourceVSphereContentLibraryRead(d, meta)
}

func resourceVSphereContentLibraryDelete(d *schema.ResourceData, meta interface{}) error {
	c := meta.(*VSphereClient).restClient
	lib, err := contentlibrary.FromName(c, d.Get("name").(string))
	if err != nil {
		return err
	}
	return contentlibrary.DeleteLibrary(c, lib)
}
