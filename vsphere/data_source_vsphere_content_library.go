package vsphere

import (
	"context"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/terraform-providers/terraform-provider-vsphere/vsphere/internal/helper/provider"
	"github.com/vmware/govmomi/vapi/library"
	"github.com/vmware/govmomi/vapi/vcenter"
	"log"
)

func dataSourceVSphereContentLibrary() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceVSphereContentLibraryRead,
		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The name or absolute path to the cluster.",
			},
			"template_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dataSourceVSphereContentLibraryRead(d *schema.ResourceData, meta interface{}) error {

	find := library.Find{
		Name: "VirtuallyGhetto",
	}
	w, err := clm.FindLibrary(ctx, find)
	log.Printf("[DEBUG] BILLLLLLLLLLLLLLLLLLLLLLLLLL--------------------- %v:%v", w, err)
	fi := library.FindItem{
		LibraryID: w[0],
		Name:      "Nested_ESXi6.0u3_Appliance_Template_v1.0",
	}
	ci, err := clm.FindLibraryItems(ctx, fi)
	item, _ := clm.GetLibraryItem(ctx, ci[0])
	log.Printf("[DEBUG] BILLLLLLLLLLLLLLLLLLLLLLLLLL--------------------- %v:%v", item, err)
	vcenter.Deploy
m := vcenter.NewManager(meta.(*VSphereClient).restClient)
vcenter.Deploy{
	DeploymentSpec: vcenter.DeploymentSpec{},
	Target:         vcenter.Target{},
}
clm.GetItem

	return nil
}

func resourceVSphereContentLibraryClient(meta interface{}) *library.Manager {
	return library.NewManager(meta.(*VSphereClient).restClient)
}

func IsLibraryItem(meta interface{}, id string) bool {
	clm := resourceVSphereContentLibraryClient(meta)
	ctx, cancel := context.WithTimeout(context.Background(), provider.DefaultAPITimeout)
	defer cancel()
	if item, _ := clm.GetLibraryItem(ctx, id); item != nil{
		return true
	}
	return false

}
