<template>
  <el-row class="app-container">
    <el-col :span="8">
      <el-button icon="el-icon-plus" @click="openEditServerDialog()">添加新的redis服务</el-button>
      <el-tree
        ref="tree"
        v-loading="treeDataLoading"
        :data="treeData"
        lazy
        :props="{isLeaf: 'isLeaf'}"
        node-key="id"
        :load="loadRedisKeyTreeNodes"
        @node-click="onTreeNodeClick"
      >
        <div slot-scope="{node, data}" class="custom-tree-node">
          <span>{{node.label}}</span>
          <span class="toolbar">
            <el-button v-if="node.level !== 1" size="mini" icon="el-icon-plus" @click.stop="clickAdd(node, data)"></el-button>
            <!-- 仅在服务、树叶上显示编辑 -->
            <el-button v-if="node.level === 1 || node.isLeaf" size="mini" icon="el-icon-edit" @click.stop="clickEdit(node, data)"></el-button>
            <el-button size="mini" icon="el-icon-refresh" @click.stop="clickReload(node, data)"></el-button>
            <!-- DB不能删除 -->
            <el-button v-if="node.level !== 2" size="mini" icon="el-icon-delete" @click.stop="clickDel(node, data)"></el-button>
          </span>
        </div>
      </el-tree>
    </el-col>

    <el-col :span="16">
      <el-tabs v-model="currentTabId" type="card" closable @tab-remove="removeTab">
        <el-tab-pane v-for="item in valuePanes" :key="item.id" :label="item.key" :name="item.id">
          <value-view :key-info="item" />
        </el-tab-pane>
      </el-tabs>
    </el-col>

    <edit-server-dialog ref="editServerDialog" @reload="loadServerOpts"/>
  </el-row>
</template>
<script>
import { listServer } from '@/api/server'
import { getDbCount, keyTreeNodes } from '@/api/redis'
import EditServerDialog from './editServerDialog'
import ValueView from './valueView'

export default {
  components: { EditServerDialog, ValueView },
  data() {
    return {
      treeData: [],
      treeDataLoading: true,
      keyContextmenuVisible: false,
      valuePanes: [],
      currentTabId: '',
      paneId: 0
    }
  },
  watch: {
  },
  mounted() {
  },
  methods: {
    async loadRedisKeyTreeNodes(node, resolve) {
      console.log('loadRedisKeyTreeNodes', node)
      const serverId = node.data.serverId
      const db = node.data.db
      const prefix = node.data.prefix
      switch (node.level) {
        case 0:
          listServer().then(data => {
            resolve(data.map(s => { return { id: 'server_' + s.id, label: s.name, isLeaf: false, serverId: s.id } }))
            this.treeDataLoading = false
          })
          break
        case 1:
          getDbCount({ serverId: serverId }).then(data => {
            const size = data.number * 1
            const nodes = []
            for (let i = 0; i < size; i++) {
              nodes.push({ id: 'DB' + i, label: 'DB' + i, isLeaf: false, serverId: serverId, db: i })
            }
            resolve(nodes)
          })
          break
        default:
          keyTreeNodes({ serverId: serverId, db: db, prefix: prefix }).then(data => {
            resolve(data.children.map(n => {
              n.serverId = serverId
              n.db = db
              n.children = []
              n.prefix = prefix ? prefix + ':' + n.label : n.label
              n.id = 'key_' + n.prefix
              if (n.isLeaf) {
                n.label = n.prefix
                n.children = null
              }
              return n
            }))
          })
      }
    },
    openEditServerDialog() {
      this.$refs['editServerDialog'].open()
    },
    loadServerOpts() {

    },
    clickAdd(node, data) {
      console.log('clickAdd', node, data)
    },
    clickEdit(node, data) {
      console.log('clickEdit', node, data)
    },
    clickDel(node, data) {
      console.log('clickDel', node, data)
    },
    clickReload(node, data) {
      node.loaded = false
      node.expand()
    },
    async onTreeNodeClick(data, node, tree) {
      if (!node.isLeaf) {
        return
      }
      const pane = {
        serverId: data.serverId,
        db: data.db,
        key: data.prefix
      }
      let index = this.valuePanes.findIndex(p => p.serverId === pane.serverId && p.db === pane.db && p.key === pane.key)
      if (index < 0) {
        pane.id = '' + this.paneId++
        this.valuePanes.push(pane)
        index = this.valuePanes.length - 1
      }
      this.currentTabId = this.valuePanes[index].id
    },
    removeTab(id) {
      const index = this.valuePanes.findIndex(p => p.id === id)
      const findPane = this.valuePanes[index]
      this.valuePanes.splice(index, 1)
      if (this.currentTabId === findPane.id) {
        if (this.valuePanes.length === 0) {
          this.currentTabId = ''
        } else if (index >= this.valuePanes.length) {
          this.currentTabId = this.valuePanes[this.valuePanes.length - 1].id
        } else {
          this.currentTabId = this.valuePanes[index].id
        }
      }
    }
  }
}
</script>

<style scoped>
  .custom-tree-node {
    flex: 1;
    display: flex;
    align-items: center;
    justify-content: space-between;
    font-size: 14px;
    padding-right: 8px;
  }
  .custom-tree-node:hover .toolbar {
    display: block;
  }
  .toolbar {
    float: right;
    display: none;
  }
  .el-button--mini {
    padding: 1px 1px;
  }
</style>
