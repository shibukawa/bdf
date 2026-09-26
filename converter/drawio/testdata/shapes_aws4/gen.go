//go:build ignore

// Gen writes the test diagrams of the AWS shapes (mxgraph.aws4.productIcon,
// resourceIcon, group, groupCenter and group2) of this directory: the
// icons in the category colors of the palette, gradients, colors, sizes,
// opacities, rotations, directions and flips, and the groups of the
// palette with their variants (grIcon, grIconSize, grStroke, strokeOpacity).
// Labels name the styles (draw.io and the converter both draw them, which
// helps to compare an export shape by shape).
//
//	go run ./converter/drawio/testdata/shapes_aws4/gen.go converter/drawio/testdata/shapes_aws4
//	/Applications/draw.io.app/Contents/MacOS/draw.io -x -f png -s 2 -o ref.png converter/drawio/testdata/shapes_aws4/aws4_icons.drawio
//	go run ./cmd/bdf generate -q converter/drawio/testdata/shapes_aws4/aws4_icons.drawio out.bdf
//	node test/render.mjs out.bdf out/ 2.6667
package main

import (
	"fmt"
	"html"
	"os"
	"path/filepath"
	"strings"
)

type doc struct {
	b  strings.Builder
	id int
}

func (d *doc) next() string { d.id++; return fmt.Sprint("c", d.id) }

func (d *doc) cell(style string, x, y, w, h float64, value string) {
	fmt.Fprintf(&d.b, `        <mxCell id="%s" value="%s" style="%s" vertex="1" parent="1">
          <mxGeometry x="%g" y="%g" width="%g" height="%g" as="geometry"/>
        </mxCell>
`, d.next(), html.EscapeString(value), html.EscapeString(style), x, y, w, h)
}

// text writes a small label centered in the box x, y, w, 24.
func (d *doc) text(x, y, w float64, s string) {
	d.cell("text;html=1;align=center;verticalAlign=top;fontSize=8;spacing=0;whiteSpace=wrap;", x, y, w, 24, s)
}

// row places a shape for each style in a row of cells cw wide, w by h,
// with the label below; it returns the y below the row.
func (d *doc) row(y0, cw, w, h float64, styles, labels []string) float64 {
	for i, s := range styles {
		x := 20 + float64(i)*cw
		d.cell(s, x+(cw-w)/2, y0, w, h, "")
		d.text(x, y0+h+4, cw, labels[i])
	}
	return y0 + h + 40
}

func (d *doc) write(path string) {
	out := `<mxfile host="gen">
  <diagram id="p1" name="Page-1">
    <mxGraphModel grid="0" gridSize="10" guides="1" tooltips="1" connect="1" arrows="1" fold="1" page="0" pageScale="1" pageWidth="827" pageHeight="1169" math="0" shadow="0">
      <root>
        <mxCell id="0"/>
        <mxCell id="1" parent="0"/>
` + d.b.String() + `      </root>
    </mxGraphModel>
  </diagram>
</mxfile>
`
	if err := os.WriteFile(path, []byte(out), 0o644); err != nil {
		panic(err)
	}
}

// variants applies each variant to base and labels the shapes with it.
func variants(base string, vs []string) (styles, labels []string) {
	for _, v := range vs {
		styles = append(styles, base+v)
		labels = append(labels, v)
	}
	return styles, labels
}

const (
	// the styles of the palette (Sidebar-AWS4.js and Sidebar-AWS4b.js)
	resource = "sketch=0;outlineConnect=0;fontColor=#232F3E;strokeColor=#ffffff;dashed=0;verticalLabelPosition=bottom;verticalAlign=top;align=center;html=1;fontSize=12;fontStyle=0;aspect=fixed;shape=mxgraph.aws4.resourceIcon;"
	product  = "sketch=0;outlineConnect=0;fontColor=#232F3E;gradientColor=none;strokeColor=#ffffff;fillColor=#232F3E;dashed=0;verticalLabelPosition=middle;verticalAlign=bottom;align=center;html=1;whiteSpace=wrap;fontSize=10;fontStyle=1;spacing=3;shape=mxgraph.aws4.productIcon;"
	group    = "points=[[0,0],[0.25,0],[0.5,0],[0.75,0],[1,0],[1,0.25],[1,0.5],[1,0.75],[1,1],[0.75,1],[0.5,1],[0.25,1],[0,1],[0,0.75],[0,0.5],[0,0.25]];outlineConnect=0;gradientColor=none;html=1;whiteSpace=wrap;fontSize=12;fontStyle=0;container=1;pointerEvents=0;collapsible=0;recursiveResize=0;"
)

func main() {
	dir := os.Args[1]

	// resourceIcon and productIcon
	d := &doc{}
	y := 0.0
	var styles, labels []string
	for _, c := range []struct{ fill, icon string }{
		{"fillColor=#ED7100;", "ec2"}, {"fillColor=#01A88D;", "sagemaker"}, {"fillColor=#DD344C;", "guardduty"},
		{"fillColor=#8C4FFF;", "athena"}, {"fillColor=#7AA116;", "s3"}, {"fillColor=#C925D1;", "dynamodb"},
		{"fillColor=#E7157B;", "cloudwatch_2"}, {"fillColor=#232F3E;", "general"},
		{"gradientColor=none;gradientDirection=north;fillColor=#1E262E;", "marketplace"},
		{"gradientColor=#4D72F3;gradientDirection=north;fillColor=#3334B9;", "iot_core"},
		{"gradientColor=#F78E04;gradientDirection=north;fillColor=#D05C17;", "lambda"},
		{"gradientColor=#F34482;gradientDirection=north;fillColor=#BC1356;", "sns"},
	} {
		styles = append(styles, resource+c.fill+"resIcon=mxgraph.aws4."+c.icon+";")
		labels = append(labels, c.fill+c.icon)
	}
	y = d.row(y, 80, 60, 60, styles, labels)

	styles, labels = variants(resource+"fillColor=#ED7100;resIcon=mxgraph.aws4.lambda;", []string{
		"", "rotation=30;", "direction=south;", "direction=north;", "direction=west;", "flipH=1;",
		"flipV=1;", "flipH=1;flipV=1;", "direction=south;flipH=1;", "direction=north;flipV=1;rotation=20;",
	})
	y = d.row(y, 80, 60, 60, styles, labels)

	styles, labels = variants(resource+"resIcon=mxgraph.aws4.users;", []string{
		"strokeColor=#232F3E;fillColor=#ffffff;", "strokeColor=none;fillColor=#FF9900;", "fillColor=none;strokeColor=#232F3E;",
		"fillColor=#ED7100;opacity=50;", "fillColor=#ED7100;fillOpacity=40;", "fillColor=#ED7100;strokeOpacity=30;",
		"fillColor=#ED7100;shadow=1;", "fillColor=#ED7100;pointerEvents=1;", "fillColor=#ED7100;strokeWidth=4;dashed=1;",
		"fillColor=#ED7100;rounded=1;glass=1;",
	})
	y = d.row(y, 80, 60, 60, styles, labels)

	styles, labels = variants(resource+"fillColor=#8C4FFF;gradientColor=#F34482;resIcon=mxgraph.aws4.api_gateway;", []string{
		"gradientDirection=south;", "gradientDirection=east;", "gradientDirection=west;", "gradientDirection=radial;",
		"resIcon=mxgraph.aws4.piop;", "resIcon=;", "resIcon=mxgraph.aws4.resourceIcon;", "resIcon=mxgraph.flowchart.decision;",
	})
	y = d.row(y, 80, 60, 60, styles, labels)

	// sizes: the stencils keep their aspect in the middle 80%
	for i, sz := range [][2]float64{{100, 60}, {30, 60}, {24, 24}, {120, 120}} {
		x := 20 + float64(i)*140
		d.cell(resource+"fillColor=#C925D1;resIcon=mxgraph.aws4.rds;", x, y, sz[0], sz[1], fmt.Sprintf("%gx%g", sz[0], sz[1]))
	}
	y += 170

	// productIcon (Sidebar-AWS4b.js): the label in the frame below the square
	styles, labels = variants(product+"prIcon=mxgraph.aws4.athena;", []string{
		"", "fillColor=#D05C17;gradientColor=#F78E04;gradientDirection=north;", "fillColor=none;strokeColor=#232F3E;",
		"strokeColor=none;", "opacity=50;", "fillOpacity=40;", "rotation=30;", "flipH=1;", "direction=south;",
		"prIcon=mxgraph.aws4.piop;",
	})
	for i, s := range styles {
		d.cell(s, 20+float64(i)*100+10, y, 80, 110, "Amazon Athena")
		d.text(20+float64(i)*100, y+114, 100, labels[i])
	}
	y += 160
	for i, sz := range [][2]float64{{80, 100}, {60, 60}, {100, 80}, {40, 120}} {
		d.cell(product+"prIcon=mxgraph.aws4.ec2;fillColor=#D05C17;", 20+float64(i)*140, y, sz[0], sz[1], "EC2")
		d.text(20+float64(i)*140, y+sz[1]+26, 100, fmt.Sprintf("%gx%g", sz[0], sz[1]))
	}
	d.write(filepath.Join(dir, "aws4_icons.drawio"))

	// the groups of the palette (Sidebar-AWS4.js and Sidebar-AWS4b.js)
	d = &doc{}
	pal := []struct{ style, label string }{
		{"group;grIcon=mxgraph.aws4.group_aws_cloud_alt;strokeColor=#232F3E;fillColor=none;verticalAlign=top;align=left;spacingLeft=30;fontColor=#232F3E;dashed=0;", "AWS Cloud"},
		{"group;grIcon=mxgraph.aws4.group_aws_cloud;strokeColor=#232F3E;fillColor=none;verticalAlign=top;align=left;spacingLeft=30;fontColor=#232F3E;dashed=0;", "AWS Cloud"},
		{"group;grIcon=mxgraph.aws4.group_region;strokeColor=#00A4A6;fillColor=none;verticalAlign=top;align=left;spacingLeft=30;fontColor=#147EBA;dashed=1;", "Region"},
		{"groupCenter;grIcon=mxgraph.aws4.group_auto_scaling_group;grStroke=1;strokeColor=#D86613;fillColor=none;verticalAlign=top;align=center;fontColor=#D86613;dashed=1;spacingTop=25;", "Auto Scaling group"},
		{"group;grIcon=mxgraph.aws4.group_vpc2;strokeColor=#8C4FFF;fillColor=none;verticalAlign=top;align=left;spacingLeft=30;fontColor=#AAB7B8;dashed=0;", "VPC"},
		{"group;grIcon=mxgraph.aws4.group_security_group;grStroke=0;strokeColor=#00A4A6;fillColor=#E6F6F7;verticalAlign=top;align=left;spacingLeft=30;fontColor=#147EBA;dashed=0;", "Private subnet"},
		{"group;grIcon=mxgraph.aws4.group_security_group;grStroke=0;strokeColor=#7AA116;fillColor=#F2F6E8;verticalAlign=top;align=left;spacingLeft=30;fontColor=#248814;dashed=0;", "Public subnet"},
		{"group;grIcon=mxgraph.aws4.group_on_premise;strokeColor=#7D8998;fillColor=none;verticalAlign=top;align=left;spacingLeft=30;fontColor=#5A6C86;dashed=0;", "Server contents"},
		{"group;grIcon=mxgraph.aws4.group_corporate_data_center;strokeColor=#7D8998;fillColor=none;verticalAlign=top;align=left;spacingLeft=30;fontColor=#5A6C86;dashed=0;", "Corporate data center"},
		{"group;grIcon=mxgraph.aws4.group_elastic_beanstalk;strokeColor=#D86613;fillColor=none;verticalAlign=top;align=left;spacingLeft=30;fontColor=#D86613;dashed=0;", "Elastic Beanstalk container"},
		{"group;grIcon=mxgraph.aws4.group_ec2_instance_contents;strokeColor=#D86613;fillColor=none;verticalAlign=top;align=left;spacingLeft=30;fontColor=#D86613;dashed=0;", "EC2 instance contents"},
		{"group;grIcon=mxgraph.aws4.group_spot_fleet;strokeColor=#D86613;fillColor=none;verticalAlign=top;align=left;spacingLeft=30;fontColor=#D86613;dashed=0;", "Spot Fleet"},
		{"group;grIcon=mxgraph.aws4.group_aws_step_functions_workflow;strokeColor=#CD2264;fillColor=none;verticalAlign=top;align=left;spacingLeft=30;fontColor=#CD2264;dashed=0;", "AWS Step Functions workflow"},
		{"group;grIcon=mxgraph.aws4.group_account;strokeColor=#CD2264;fillColor=none;verticalAlign=top;align=left;spacingLeft=30;fontColor=#CD2264;dashed=0;", "AWS Account"},
		{"group;grIcon=mxgraph.aws4.group_iot_greengrass_deployment;strokeColor=#7AA116;fillColor=none;verticalAlign=top;align=left;spacingLeft=30;fontColor=#3F8624;dashed=0;", "AWS IoT Greengrass Deployment"},
		{"group;grIcon=mxgraph.aws4.group_iot_greengrass;strokeColor=#7AA116;fillColor=none;verticalAlign=top;align=left;spacingLeft=30;fontColor=#3F8624;dashed=0;", "AWS IoT Greengrass"},
		{"group;grIcon=mxgraph.aws4.group_availability_zone;strokeColor=#545B64;fillColor=none;verticalAlign=top;align=left;spacingLeft=30;fontColor=#545B64;dashed=1;", "Availability zone"},
		{"group;grIcon=mxgraph.aws4.group_subnet;strokeColor=#879196;fillColor=none;verticalAlign=top;align=left;spacingLeft=30;fontColor=#879196;dashed=0;", "Subnet"},
		{"groupCenter;grIcon=mxgraph.aws4.group_auto_scaling_group;grStroke=0;strokeColor=#879196;fillColor=#ECEFEF;verticalAlign=top;align=center;fontColor=#879196;dashed=0;spacingTop=25;", "Auto Scaling Group"},
		{"groupCenter;grIcon=mxgraph.aws4.group_elastic_load_balancing;grStroke=1;strokeColor=#007DBC;fillColor=none;verticalAlign=top;align=center;fontColor=#007DBC;dashed=0;spacingTop=25;", "Elastic Load Balancing"},
	}
	for i, p := range pal {
		d.cell(group+"shape=mxgraph.aws4."+p.style, 20+float64(i%7)*150, float64(i/7)*150, 130, 130, p.label)
	}
	y = 3*150 + 10

	// group and groupCenter variants
	for i, v := range []string{
		"grIconSize=40;", "grIconSize=15;", "grIconSize=abc;", "grStroke=2;fillColor=#E6F6F7;", "strokeOpacity=40;",
		"strokeColor=none;", "strokeWidth=3;", "shadow=1;", "fillColor=#E6F6F7;gradientColor=#7AA116;",
		"rotation=20;", "flipH=1;", "flipV=1;", "direction=south;", "direction=west;",
	} {
		for j, shape := range []string{"group", "groupCenter"} {
			style := group + "shape=mxgraph.aws4." + shape + ";grIcon=mxgraph.aws4.group_aws_cloud_alt;strokeColor=#232F3E;fillColor=none;verticalAlign=top;align=left;spacingLeft=30;fontColor=#232F3E;" + v
			x := 20 + float64(i%7)*150
			yy := y + float64(i/7)*220 + float64(j)*100
			d.cell(style, x, yy, 130, 80, shape)
			if j == 1 {
				d.text(x, yy+84, 130, v)
			}
		}
	}
	y += 2*220 + 10

	// group2 (in no palette) and small groups
	for i, v := range []string{
		"", "grIconSize=40;", "strokeOpacity=40;", "strokeColor=#ED7100;fillColor=#FFF2E6;", "grIcon=;",
		"rotation=20;", "flipH=1;",
	} {
		style := group + "shape=mxgraph.aws4.group2;grIcon=mxgraph.aws4.group_aws_cloud_alt;strokeColor=#232F3E;fillColor=none;verticalAlign=top;align=left;spacingLeft=30;fontColor=#232F3E;" + v
		x := 20 + float64(i)*150
		d.cell(style, x, y, 130, 80, "group2")
		d.text(x, y+84, 130, v)
	}
	y += 120
	for i, shape := range []string{"group", "groupCenter", "group2"} {
		style := group + "shape=mxgraph.aws4." + shape + ";grIcon=mxgraph.aws4.group_region;strokeColor=#00A4A6;fillColor=none;dashed=1;"
		d.cell(style, 20+float64(i)*60, y, 20, 20, "")
		d.cell(style, 20+float64(i)*60, y+40, 40, 10, "")
	}
	d.write(filepath.Join(dir, "aws4_groups.drawio"))
}
