// Code generated by templ - DO NOT EDIT.

// templ: version: v0.2.543
package web

//lint:file-ignore SA4006 This context is only used if a nested component is present.

import "github.com/a-h/templ"
import "context"
import "io"
import "bytes"

func gameTable() templ.ComponentScript {
	return templ.ComponentScript{
		Name: `__templ_gameTable_86b9`,
		Function: `function __templ_gameTable_86b9(){const format = (d) => {
		return '<dl>' +
				'<dt>Title:</dt>' +
				'<dd>' +
				d.title +
				'</dd>' +
			'</dl>';
	}

	var table = new DataTable('#games-table', {
		"processing": true,
        "serverSide": true,
        "ajax": "/games/table",
        "columns": [
			{
				className: "dt-control",
				orderable: false,
				data: null,
				defaultContent: "",
			},
            { 
				data: "title"
			},
            {
				data: "category"
			},
            {
				data: "platforms",
				render: function(data, type, row) {
					// If platforms is an array, join it into a string
					return data.join(", ");
            	}
			},
            {
				data: "downloaded", 
				render: function(data, type, row) {
					return ` + "`" + `<span class="inline-flex flex-shrink-0 items-center rounded-full bg-${data ? "green" : "red"}-50 px-1.5 py-0.5 text-xs font-medium text-${data ? "green" : "red"}-700 ring-1 ring-inset ring-${data ? "green" : "red"}-600/20">${data ? "downloaded" : "no backup"}</span>` + "`" + `
            	}
			}
        ],
		///
		responsive: true,
		select: true,
		layout: {
			topStart: {
				buttons: ['pageLength', 'selectAll', 'selectNone'] // Add 'selectNone' button
			}
		},
		buttons: [
			{
				extend: 'selected',
				text: 'Count selected rows',
				action: function (e, dt, button, config) {
					alert(
						dt.rows({ selected: true }).indexes().length + ' row(s) selected'
					);
				}
			}
		],
		// "initComplete": function(settings, json) {
		// 	// Target the DataTable wrapper
		// 	// $(this).closest('.dataTables_wrapper').addClass('myCustomWrapperClass');
			
		// 	// Target the table and its components directly
		// 	var api = this.api();

		// 	$(api.table()).addClass('leading-normal');
		// 	// $(api.table().header()).addClass('myCustomHeaderClass');
		// 	// $(api.table().body()).addClass('myCustomBodyClass');
			
		// 	// Apply classes to each row
		// 	// api.rows().every(function() {
		// 	// 	$(this.node()).addClass('myCustomRowClass');
		// 	// });

		// 	// Apply classes to a specific column, e.g., the first column
		// 	// $(api.columns(0).header()).addClass('myCustomColumnHeaderClass');
		// 	// api.columns(0).nodes().flatten().to$().addClass('myCustomColumnClass');
			
		// 	// Additional customizations can be added here
		// },
		// pagingType: "scrolling"
	});

	new $.fn.dataTable.FixedHeader(table);

	// Add event listener for opening and closing details
	table.on('click', 'td.dt-control', function (e) {
		let tr = e.target.closest('tr');
		let row = table.row(tr);
	
		if (row.child.isShown()) {
			// This row is already open - close it
			row.child.hide();
		}
		else {
			// Open this row
			row.child(format(row.data())).show();
		}
	});
}`,
		Call:       templ.SafeScript(`__templ_gameTable_86b9`),
		CallInline: templ.SafeScriptInline(`__templ_gameTable_86b9`),
	}
}

func PageGames(loading bool) templ.Component {
	return templ.ComponentFunc(func(ctx context.Context, templ_7745c5c3_W io.Writer) (templ_7745c5c3_Err error) {
		templ_7745c5c3_Buffer, templ_7745c5c3_IsBuffer := templ_7745c5c3_W.(*bytes.Buffer)
		if !templ_7745c5c3_IsBuffer {
			templ_7745c5c3_Buffer = templ.GetBuffer()
			defer templ.ReleaseBuffer(templ_7745c5c3_Buffer)
		}
		ctx = templ.InitializeContext(ctx)
		templ_7745c5c3_Var1 := templ.GetChildren(ctx)
		if templ_7745c5c3_Var1 == nil {
			templ_7745c5c3_Var1 = templ.NopComponent
		}
		ctx = templ.ClearChildren(ctx)
		templ_7745c5c3_Var2 := templ.ComponentFunc(func(ctx context.Context, templ_7745c5c3_W io.Writer) (templ_7745c5c3_Err error) {
			templ_7745c5c3_Buffer, templ_7745c5c3_IsBuffer := templ_7745c5c3_W.(*bytes.Buffer)
			if !templ_7745c5c3_IsBuffer {
				templ_7745c5c3_Buffer = templ.GetBuffer()
				defer templ.ReleaseBuffer(templ_7745c5c3_Buffer)
			}
			_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("<div class=\"container\">")
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
			if loading {
				_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("<div class=\"alert alert-info\" role=\"alert\">Downloading games...</div>")
				if templ_7745c5c3_Err != nil {
					return templ_7745c5c3_Err
				}
			}
			_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("<div class=\"px-4 sm:px-6 lg:px-8\"><div class=\"mt-8 flow-root\"><div class=\"-mx-4 -my-2 overflow-x-auto sm:-mx-6 lg:-mx-8\"><div class=\"inline-block min-w-full py-2 align-middle sm:px-6 lg:px-8\"><div class=\"relative\"><table id=\"games-table\" class=\"display compact\"><thead><tr><th>Title</th><th>Title</th><th>Category</th><th>Platforms</th><th>Downloaded</th></tr></thead> <tbody></tbody></table></div></div></div></div></div>")
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
			templ_7745c5c3_Err = gameTable().Render(ctx, templ_7745c5c3_Buffer)
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
			_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("</div>")
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
			if !templ_7745c5c3_IsBuffer {
				_, templ_7745c5c3_Err = io.Copy(templ_7745c5c3_W, templ_7745c5c3_Buffer)
			}
			return templ_7745c5c3_Err
		})
		templ_7745c5c3_Err = BasePage(
			"Gogogo - Games",
			NavigationBar{
				Home: true,
			},
		).Render(templ.WithChildren(ctx, templ_7745c5c3_Var2), templ_7745c5c3_Buffer)
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		if !templ_7745c5c3_IsBuffer {
			_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteTo(templ_7745c5c3_W)
		}
		return templ_7745c5c3_Err
	})
}
