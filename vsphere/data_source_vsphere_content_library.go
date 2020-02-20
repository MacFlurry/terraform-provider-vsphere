package vsphere

import (
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/terraform-providers/terraform-provider-vsphere/vsphere/internal/helper/contentlibrary"
)

func dataSourceVSphereContentLibrary() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceVSphereContentLibraryRead,
		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The name of the content library.",
			},
		},
	}
}

func dataSourceVSphereContentLibraryRead(d *schema.ResourceData, meta interface{}) error {
	c := meta.(*VSphereClient).restClient
	lib, err := contentlibrary.FromName(c, d.Get("name").(string))
	if err != nil {
		return err
	}
	d.SetId(lib.ID)
	return nil
}
