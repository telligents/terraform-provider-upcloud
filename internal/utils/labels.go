package utils

import (
	"context"
	"fmt"
	"regexp"

	"github.com/UpCloudLtd/upcloud-go-api/v8/upcloud"
	"github.com/hashicorp/terraform-plugin-framework-validators/mapvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/mapplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	sdkv2_schema "github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

var labelKeyRegExp = regexp.MustCompile("^([a-zA-Z0-9])+([a-zA-Z0-9_-])*$")

func labelsDescription(resource string) string {
	return fmt.Sprintf("Key-value pairs to classify the %s.", resource)
}

var _ planmodifier.Map = unconfiguredAsEmpty{}

type unconfiguredAsEmpty struct{}

func (lm unconfiguredAsEmpty) Description(_ context.Context) string {
	return "use empty map, if config is null."
}

func (lm unconfiguredAsEmpty) MarkdownDescription(ctx context.Context) string {
	return lm.Description(ctx)
}

func (lm unconfiguredAsEmpty) PlanModifyMap(ctx context.Context, req planmodifier.MapRequest, resp *planmodifier.MapResponse) {
	if req.ConfigValue.IsNull() {
		labels := make(map[string]string)
		value, diags := types.MapValueFrom(ctx, types.StringType, &labels)

		resp.PlanValue = value
		resp.Diagnostics = diags
	}
}

func LabelsAttribute(resource string) schema.Attribute {
	description := labelsDescription(resource)
	return &schema.MapAttribute{
		ElementType: types.StringType,
		Computed:    true,
		Optional:    true,
		Description: description,
		PlanModifiers: []planmodifier.Map{
			unconfiguredAsEmpty{},
			mapplanmodifier.UseStateForUnknown(),
		},
		Validators: []validator.Map{
			mapvalidator.KeysAre(stringvalidator.LengthBetween(2, 32), stringvalidator.RegexMatches(labelKeyRegExp, "")),
			mapvalidator.ValueStringsAre(stringvalidator.LengthBetween(0, 255)),
		},
	}
}

func LabelsSchema(resource string) *sdkv2_schema.Schema {
	description := labelsDescription(resource)
	return &sdkv2_schema.Schema{
		Description: description,
		Type:        sdkv2_schema.TypeMap,
		Elem: &sdkv2_schema.Schema{
			Type: sdkv2_schema.TypeString,
		},
		Optional:         true,
		ValidateDiagFunc: ValidateLabelsDiagFunc,
	}
}

func LabelsMapToSlice[T any](m map[string]T) []upcloud.Label {
	var labels []upcloud.Label

	for k, v := range m {
		var value any = v
		labels = append(labels, upcloud.Label{
			Key:   k,
			Value: value.(string),
		})
	}

	return labels
}

func LabelsSliceToMap(s []upcloud.Label) map[string]string {
	labels := make(map[string]string)

	for _, l := range s {
		labels[l.Key] = l.Value
	}

	return labels
}

var ValidateLabelsDiagFunc = validation.AllDiag(
	validation.MapKeyLenBetween(2, 32),
	validation.MapKeyMatch(labelKeyRegExp, ""),
	validation.MapValueLenBetween(0, 255),
)
