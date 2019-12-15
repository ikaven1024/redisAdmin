<template>
  <el-dialog :visible.sync="visible" :close-on-click-modal="false">
    <el-form label-width="80px">
      <el-form-item label="名称">
        <el-input v-model="data.name" />
      </el-form-item>
      <el-form-item label="模式">
        <el-radio-group v-model="data.mode">
          <el-radio :label="0">单节点</el-radio>
          <el-radio :label="1">集群</el-radio>
        </el-radio-group>
      </el-form-item>
      <el-form-item label="节点">
        <div v-for="(item, index) in data.addresses" :key="index">
          <el-input v-model="data.addresses[index]" placeholder="127.0.0.1:6379" style="width: auto" />
          <el-button v-if="index > 0" icon="el-icon-delete" @click="removeNode(index)" />
        </div>
        <el-button v-show="data.mode === 1" icon="el-icon-plus" style="margin-left: 10px" @click="addNode" />
      </el-form-item>
      <el-form-item label="密码">
        <el-input v-model="data.password" placeholder="密码" />
      </el-form-item>
    </el-form>

    <div slot="footer">
      <el-button type="danger" @click="removeServer">删除</el-button>
      <el-button type="primary" @click="commit">确认</el-button>
      <el-button @click="close">关闭</el-button>
    </div>
  </el-dialog>
</template>

<script>
import { addServer, updateServer, removeServer } from '@/api/server'
export default {
  data() {
    return {
      visible: false,
      data: {
        id: '',
        name: '',
        mode: 0,
        addresses: []
      }
    }
  },
  watch: {
    'data.mode'(v) {
      if (v === 0) {
        this.data.addresses = [this.data.addresses[0]]
      }
    }
  },
  methods: {
    open(data) {
      this.data = data
        ? JSON.parse(JSON.stringify(data))
        : {
          name: '',
          mode: 0,
          password: '',
          addresses: ['']
        }
      this.visible = true
    },
    addNode() {
      this.data.addresses.push('')
    },
    removeNode(index) {
      this.data.addresses.splice(index, 1)
    },
    close() {
      this.visible = false
      this.data = {}
    },
    closeAndReload(message) {
      message && this.$message.success(message)
      this.$emit('reload')
      this.close()
    },
    commit() {
      const commitData = JSON.parse(JSON.stringify(this.data))
      if (commitData.mode === 0) {
        commitData.addresses = commitData.addresses.slice(0, 1)
      }

      const promise = commitData.id ? updateServer(commitData) : addServer(commitData)
      promise.then(() => { this.closeAndReload('提交成功') })
    },
    removeServer() {
      this.$confirm(`确定要删除服务${this.data.name}吗？`, {
        confirmButtonText: '确定',
        cancelButtonText: '取消',
        type: 'warning'
      }, '提示').then(() => {
        removeServer(this.data.id).then(() => { this.closeAndReload('删除成功') })
      }).catch(() => {})
    }
  }
}
</script>

<style lang="scss" scoped>
</style>
