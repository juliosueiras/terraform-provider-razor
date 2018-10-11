provider "razor" {}

data "razor_task" "default" {
	name = "debian"
}

data "razor_node" "default" {
	name = "node4"
}

resource "razor_repo" "default" {
	name = "test20"
	task = "${data.razor_task.default.name}"
	no_content = true
}

resource "razor_tag" "default" {
	name = "esxi2"
	rule = <<EOF
[
	"=",
	["metadata","hw_info","mac"],
	"sd"
]
EOF
}

resource "razor_policy" "default" {
	name = "test_policy"
	broker = "noop"
	repo = "${razor_repo.default.name}"
	enabled = true
	hostname = "SD"
	max_count = 1
	root_password = "SD"
	task = "${data.razor_task.default.name}"
	node_metadata {
		"2" = "X"
		"SX" = "SW"
	}
}
