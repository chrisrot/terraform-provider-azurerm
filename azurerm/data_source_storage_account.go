package azurerm

import (
	"fmt"
	"strings"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func dataSourceArmStorageAccount() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceArmStorageAccountRead,

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},

			"resource_group_name": azure.SchemaResourceGroupNameForDataSource(),

			"location": azure.SchemaLocationForDataSource(),

			"account_kind": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"account_tier": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"account_replication_type": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"access_tier": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"account_encryption_source": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"custom_domain": {
				Type:     schema.TypeList,
				Computed: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},

			"enable_blob_encryption": {
				Type:     schema.TypeBool,
				Computed: true,
			},

			"enable_file_encryption": {
				Type:     schema.TypeBool,
				Computed: true,
			},

			"enable_https_traffic_only": {
				Type:     schema.TypeBool,
				Computed: true,
			},

			"is_hns_enabled": {
				Type:     schema.TypeBool,
				Computed: true,
			},

			"primary_location": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"secondary_location": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"primary_blob_endpoint": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"primary_blob_host": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"secondary_blob_endpoint": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"secondary_blob_host": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"primary_queue_endpoint": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"primary_queue_host": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"secondary_queue_endpoint": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"secondary_queue_host": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"primary_table_endpoint": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"primary_table_host": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"secondary_table_endpoint": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"secondary_table_host": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"primary_web_endpoint": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"primary_web_host": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"secondary_web_endpoint": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"secondary_web_host": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"primary_dfs_endpoint": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"primary_dfs_host": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"secondary_dfs_endpoint": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"secondary_dfs_host": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"primary_file_endpoint": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"primary_file_host": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"secondary_file_endpoint": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"secondary_file_host": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"primary_access_key": {
				Type:      schema.TypeString,
				Computed:  true,
				Sensitive: true,
			},

			"secondary_access_key": {
				Type:      schema.TypeString,
				Computed:  true,
				Sensitive: true,
			},

			"primary_connection_string": {
				Type:      schema.TypeString,
				Computed:  true,
				Sensitive: true,
			},

			"secondary_connection_string": {
				Type:      schema.TypeString,
				Computed:  true,
				Sensitive: true,
			},

			"primary_blob_connection_string": {
				Type:      schema.TypeString,
				Computed:  true,
				Sensitive: true,
			},

			"secondary_blob_connection_string": {
				Type:      schema.TypeString,
				Computed:  true,
				Sensitive: true,
			},

			"tags": tagsForDataSourceSchema(),
		},
	}

}

func dataSourceArmStorageAccountRead(d *schema.ResourceData, meta interface{}) error {
	ctx := meta.(*ArmClient).StopContext
	client := meta.(*ArmClient).storageServiceClient
	endpointSuffix := meta.(*ArmClient).environment.StorageEndpointSuffix

	name := d.Get("name").(string)
	resourceGroup := d.Get("resource_group_name").(string)

	resp, err := client.GetProperties(ctx, resourceGroup, name)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			return fmt.Errorf("Error: Storage Account %q (Resource Group %q) was not found", name, resourceGroup)
		}
		return fmt.Errorf("Error reading the state of AzureRM Storage Account %q: %+v", name, err)
	}

	d.SetId(*resp.ID)

	keys, err := client.ListKeys(ctx, resourceGroup, name)
	if err != nil {
		return err
	}

	accessKeys := *keys.Keys
	if location := resp.Location; location != nil {
		d.Set("location", azure.NormalizeLocation(*location))
	}
	d.Set("account_kind", resp.Kind)

	if sku := resp.Sku; sku != nil {
		d.Set("account_tier", sku.Tier)
		d.Set("account_replication_type", strings.Split(string(sku.Name), "_")[1])
	}

	if props := resp.AccountProperties; props != nil {
		d.Set("access_tier", props.AccessTier)
		d.Set("enable_https_traffic_only", props.EnableHTTPSTrafficOnly)
		d.Set("is_hns_enabled", props.IsHnsEnabled)

		if customDomain := props.CustomDomain; customDomain != nil {
			if err := d.Set("custom_domain", flattenStorageAccountCustomDomain(customDomain)); err != nil {
				return fmt.Errorf("Error setting `custom_domain`: %+v", err)
			}
		}

		if encryption := props.Encryption; encryption != nil {
			if services := encryption.Services; services != nil {
				if blob := services.Blob; blob != nil {
					d.Set("enable_blob_encryption", blob.Enabled)
				}
				if file := services.File; file != nil {
					d.Set("enable_file_encryption", file.Enabled)
				}
			}
			d.Set("account_encryption_source", string(encryption.KeySource))
		}

		// Computed
		d.Set("primary_location", props.PrimaryLocation)
		d.Set("secondary_location", props.SecondaryLocation)

		if len(accessKeys) > 0 {
			pcs := fmt.Sprintf("DefaultEndpointsProtocol=https;AccountName=%s;AccountKey=%s;EndpointSuffix=%s", *resp.Name, *accessKeys[0].Value, endpointSuffix)
			d.Set("primary_connection_string", pcs)
		}

		if len(accessKeys) > 1 {
			scs := fmt.Sprintf("DefaultEndpointsProtocol=https;AccountName=%s;AccountKey=%s;EndpointSuffix=%s", *resp.Name, *accessKeys[1].Value, endpointSuffix)
			d.Set("secondary_connection_string", scs)
		}

		if err := flattenAndSetAzureRmStorageAccountPrimaryEndpoints(d, props.PrimaryEndpoints); err != nil {
			return fmt.Errorf("error setting primary endpoints and hosts for blob, queue, table and file: %+v", err)
		}

		var primaryBlobConnectStr string
		if v := props.PrimaryEndpoints; v != nil {
			primaryBlobConnectStr = getBlobConnectionString(v.Blob, resp.Name, accessKeys[0].Value)
		}
		d.Set("primary_blob_connection_string", primaryBlobConnectStr)

		if err := flattenAndSetAzureRmStorageAccountSecondaryEndpoints(d, props.SecondaryEndpoints); err != nil {
			return fmt.Errorf("error setting secondary endpoints and hosts for blob, queue, table: %+v", err)
		}

		var secondaryBlobConnectStr string
		if v := props.SecondaryEndpoints; v != nil {
			secondaryBlobConnectStr = getBlobConnectionString(v.Blob, resp.Name, accessKeys[1].Value)
		}
		d.Set("secondary_blob_connection_string", secondaryBlobConnectStr)
	}

	d.Set("primary_access_key", accessKeys[0].Value)
	d.Set("secondary_access_key", accessKeys[1].Value)

	flattenAndSetTags(d, resp.Tags)

	return nil
}
