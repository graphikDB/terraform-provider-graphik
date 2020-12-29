# terraform-provider-graphik

terraform provider for graphikDB schema operations

`git clone git@github.com:graphikDB/terraform-provider-graphik.git`


## Installation

After installation, please move the binary to a location in which it may be discovered by terraform

### Mac/Darwin

```text
curl -LJO https://github.com/graphikDB/terraform-provider-graphik/releases/download/v0.1.0/terraform-provider-graphik_darwin_amd64 && \
    mv terraform-provider-graphik_darwin_amd64 terraform-provider-graphik && \
    chmod +x terraform-provider-graphik
```

### Linux

```text
curl -LJO https://github.com/graphikDB/terraform-provider-graphik/releases/download/v0.1.0/terraform-provider-graphik_linux_amd64 && \
    mv terraform-provider-graphik_darwin_amd64 terraform-provider-graphik && \
    chmod +x terraform-provider-graphik
```

## Example - Task application

```hcl-terraform
provider "graphik" {
  # graphik provider may be automatically configured after running `graphikctl auth login`
}

# graphik_constraint.task_title_description requires tasks to have a title and description
resource "graphik_constraint" "task_title_description" {
  lifecycle {
    prevent_destroy = true
  }
  name = "task_title_description"
  gtype = "task"
  expression = "has(this.attributes.title) && has(this.attributes.description)"
  target_docs = true
  target_connections = false
}

# graphik_constraint.task_priority is 
resource "graphik_constraint" "task_priority" {
  lifecycle {
    prevent_destroy = true
  }
  name = "task_priority"
  gtype = "task"
  expression = "this.attributes.priority in ['low', 'medium', 'high']"
  target_docs = true
  target_connections = false
}

# graphik_index.low_priority is a secondary index where low priority tasks can be queried from
resource "graphik_index" "low_priority" {
  lifecycle {
    prevent_destroy = true
  }
  name = "low_priority"
  gtype = "task"
  expression = "this.attributes.priority == 'low'"
  target_docs = true
  target_connections = false
}

# graphik_index.medium_priority is a secondary index where medium priority tasks can be queried from
resource "graphik_index" "medium_priority" {
  lifecycle {
    prevent_destroy = true
  }
  name = "low_priority"
  gtype = "task"
  expression = "this.attributes.priority == 'medium'"
  target_docs = true
  target_connections = false
}

# graphik_index.high_priority is a secondary index where high priority tasks can be queried from
resource "graphik_index" "high_priority" {
  lifecycle {
    prevent_destroy = true
  }
  name = "low_priority"
  gtype = "task"
  expression = "this.attributes.priority == 'high'"
  target_docs = true
  target_connections = false
}

# graphik_trigger.updated_at adds a updated_at timestamp to doc & connection attributes any time it is changed
resource "graphik_trigger" "updated_at" {
  lifecycle {
    prevent_destroy = true
  }
  name = "updated_at"
  gtype = "*"
  expression = "true"
  trigger = "{ 'updated_at': now() }"
  target_docs = true
  target_connections = true
}

# graphik_trigger.created_at adds a created_at timestamp to doc & connection attributes when it's created
resource "graphik_trigger" "created_at" {
  lifecycle {
    prevent_destroy = true
  }
  name = "created_at"
  gtype = "*"
  expression = "!has(this.attributes.created_at)"
  trigger = "{ 'created_at': now() }"
  target_docs = true
  target_connections = true
}
```