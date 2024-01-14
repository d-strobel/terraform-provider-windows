// Code generated by terraform-plugin-framework-generator DO NOT EDIT.

package local_datasources

import (
	"context"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
)

func LocalUserDataSourceSchema(ctx context.Context) schema.Schema {
	return schema.Schema{
		Attributes: map[string]schema.Attribute{
			"account_expires": schema.StringAttribute{
				Computed:            true,
				Description:         "Retrieve the time where the local user account expires.",
				MarkdownDescription: "Retrieve the time where the local user account expires.",
			},
			"description": schema.StringAttribute{
				Computed:            true,
				Description:         "The description of the local user.",
				MarkdownDescription: "The description of the local user.",
			},
			"enabled": schema.BoolAttribute{
				Computed:            true,
				Description:         "Get the status of the local user.",
				MarkdownDescription: "Get the status of the local user.",
			},
			"full_name": schema.StringAttribute{
				Computed:            true,
				Description:         "The full name of the local user.",
				MarkdownDescription: "The full name of the local user.",
			},
			"id": schema.StringAttribute{
				Computed:            true,
				Description:         "The ID of the retrieved local user. This is the same as the SID.",
				MarkdownDescription: "The ID of the retrieved local user. This is the same as the SID.",
			},
			"last_login": schema.StringAttribute{
				Computed:            true,
				Description:         "The last login time of the local user.",
				MarkdownDescription: "The last login time of the local user.",
			},
			"name": schema.StringAttribute{
				Required:            true,
				Description:         "Define the name of the local user.",
				MarkdownDescription: "Define the name of the local user.",
			},
			"password_changeable_date": schema.StringAttribute{
				Computed:            true,
				Description:         "The password changeable date of the local user.",
				MarkdownDescription: "The password changeable date of the local user.",
			},
			"password_expires": schema.StringAttribute{
				Computed:            true,
				Description:         "The time when the password of the local user expires.",
				MarkdownDescription: "The time when the password of the local user expires.",
			},
			"password_last_set": schema.StringAttribute{
				Computed:            true,
				Description:         "The last time when the password was set for the local user.",
				MarkdownDescription: "The last time when the password was set for the local user.",
			},
			"password_required": schema.BoolAttribute{
				Computed:            true,
				Description:         "If true a password is required login with the local user.",
				MarkdownDescription: "If true a password is required login with the local user.",
			},
			"sid": schema.StringAttribute{
				Optional:            true,
				Description:         "The security ID of the local user.",
				MarkdownDescription: "The security ID of the local user.",
			},
			"user_may_change_password": schema.BoolAttribute{
				Computed:            true,
				Description:         "If true the local user can change it's password.",
				MarkdownDescription: "If true the local user can change it's password.",
			},
		},
	}
}

type LocalUserModel struct {
	AccountExpires         types.String `tfsdk:"account_expires"`
	Description            types.String `tfsdk:"description"`
	Enabled                types.Bool   `tfsdk:"enabled"`
	FullName               types.String `tfsdk:"full_name"`
	Id                     types.String `tfsdk:"id"`
	LastLogin              types.String `tfsdk:"last_login"`
	Name                   types.String `tfsdk:"name"`
	PasswordChangeableDate types.String `tfsdk:"password_changeable_date"`
	PasswordExpires        types.String `tfsdk:"password_expires"`
	PasswordLastSet        types.String `tfsdk:"password_last_set"`
	PasswordRequired       types.Bool   `tfsdk:"password_required"`
	Sid                    types.String `tfsdk:"sid"`
	UserMayChangePassword  types.Bool   `tfsdk:"user_may_change_password"`
}