package openstack

import (
	"fmt"
	log "github.com/sourcegraph-ce/logrus"

	"github.com/gophercloud/gophercloud/openstack/blockstorage/v2/volumes"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func dataSourceBlockStorageVolumeV2() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceBlockStorageVolumeV2Read,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"name": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"status": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"metadata": {
				Type:     schema.TypeMap,
				Optional: true,
				Computed: true,
			},

			// Computed values
			"bootable": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"volume_type": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"size": {
				Type:     schema.TypeInt,
				Computed: true,
			},

			"source_volume_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dataSourceBlockStorageVolumeV2Read(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	client, err := config.BlockStorageV2Client(GetRegion(d, config))
	if err != nil {
		return fmt.Errorf("Error creating OpenStack block storage client: %s", err)
	}

	listOpts := volumes.ListOpts{
		Metadata: expandToMapStringString(d.Get("metadata").(map[string]interface{})),
		Name:     d.Get("name").(string),
		Status:   d.Get("status").(string),
	}

	allPages, err := volumes.List(client, listOpts).AllPages()
	if err != nil {
		return fmt.Errorf("Unable to query openstack_blockstorage_volume_v2: %s", err)
	}

	allVolumes, err := volumes.ExtractVolumes(allPages)
	if err != nil {
		return fmt.Errorf("Unable to retrieve openstack_blockstorage_volume_v2: %s", err)
	}

	if len(allVolumes) > 1 {
		return fmt.Errorf("Your openstack_blockstorage_volume_v2 query returned multiple results.")
	}

	if len(allVolumes) < 1 {
		return fmt.Errorf("Your openstack_blockstorage_volume_v2 query returned no results.")
	}

	return dataSourceBlockStorageVolumeV2Attributes(d, allVolumes[0])
}

func dataSourceBlockStorageVolumeV2Attributes(d *schema.ResourceData, volume volumes.Volume) error {
	d.SetId(volume.ID)
	d.Set("name", volume.Name)
	d.Set("status", volume.Status)
	d.Set("bootable", volume.Bootable)
	d.Set("volume_type", volume.VolumeType)
	d.Set("size", volume.Size)
	d.Set("source_volume_id", volume.SourceVolID)

	if err := d.Set("metadata", volume.Metadata); err != nil {
		log.Printf("[DEBUG] Unable to set metadata for volume %s: %s", volume.ID, err)
	}

	return nil
}
