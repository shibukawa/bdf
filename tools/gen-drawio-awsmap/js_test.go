package main

import "testing"

func TestReadPalettes(t *testing.T) {
	src := `(function(){
	Sidebar.prototype.addX = function()
	{
		var s = 'dashed=0;html=1;' + mxConstants.STYLE_SHAPE + "=mxgraph.aws2.";
		var gn = 'mxgraph.aws.analytics';
		var dt = 'aws analytics ';
		var k = 1.5; // scale
		this.addPaletteFunctions('aws2Analytics', 'AWS / Analytics', false,
		[
		 this.createVertexTemplateEntry(s + 'analytics.data_pipeline;strokeColor=none;',
				 k * 60, 72, '', 'Data Pipeline', null, null, this.getTagsForStencil(gn, 'data pipeline', dt).join(' ')),
		 this.addEntry(dt + 'cloud', function()
			{
				var bg1 = new mxCell('', new mxGeometry(0, 30, 200, 200), s + 'rrect;fillColor=none;');
				bg1.vertex = true;
			   	return sb.createVertexTemplateFromCells([bg1], 200, 230, 'AWS Cloud');
			})
		]);
	};})();`
	es, err := readPalettes(src)
	if err != nil {
		t.Fatal(err)
	}
	if len(es) != 2 || es[0].w != 90 || es[0].title != "Data Pipeline" || es[0].tags != "data pipeline" ||
		es[0].style != "dashed=0;html=1;shape=mxgraph.aws2.analytics.data_pipeline;strokeColor=none;" ||
		es[1].title != "AWS Cloud" || es[1].w != 200 || es[1].palette != "aws2Analytics" ||
		es[1].style != "dashed=0;html=1;shape=mxgraph.aws2.rrect;fillColor=none;" {
		t.Errorf("got %+v", es)
	}
}
