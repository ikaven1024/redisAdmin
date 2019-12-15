<template>
  <el-form v-loading="loading">
    <el-form-item>
      <el-input v-model="keyInfo.key">
        <template slot="prepend">{{ data.type.toUpperCase() }}</template>
        <template slot="append">
          TTL: {{ data.ttl }}
        </template>
      </el-input>
    </el-form-item>

    <el-form-item style="margin-bottom: 0">
      <el-button icon="el-icon-edit">重命名</el-button>
      <el-button icon="el-icon-edit">设置TTL</el-button>
      <el-button icon="el-icon-refresh" @click="reload">刷新</el-button>
      <el-button icon="el-icon-delete" type="danger">删除</el-button>

      <div v-if="data.total" style="float: right">
        <div v-if="data.type === 'list'">
          <el-button :disabled="query.pageNo === 1" @click="loadListPrev">上一页</el-button>
          <span>共 {{ query.pageNo }} 页</span>
          <el-button :disabled="query.pageNo * query.pageSize >= data.total" @click="loadListNext">下一页</el-button>
          <span>共 {{ data.total }} 条</span>
        </div>
        <div v-else>
          <span>共 {{ data.total }} 条</span>
          <el-button v-if="data.total" :disabled="!query.cursor" @click="loadNext">下一页</el-button>
        </div>
      </div>
    </el-form-item>

    <el-form-item>
      <el-table :data="data.value">
        <el-table-column v-if="data.type === 'hash'" label="key" prop="key" />
        <el-table-column label="value" prop="value" />
        <el-table-column v-if="data.type === 'zset'" label="score" prop="score" />
        <el-table-column>
          <template scope="scope">
            <el-button type="text" @click="openValueEditDialog(scope.row)">查看</el-button>
            <!--<el-button type="text">删除</el-button>-->
          </template>
        </el-table-column>
      </el-table>
    </el-form-item>

    <value-edit-dialog ref="valueEditDialog" />
  </el-form>
</template>

<script>
import { keySummary, getValue } from '@/api/redis'
import ValueEditDialog from './valueEditDialog'

export default {
  name: 'ValueView',
  components: { ValueEditDialog },
  props: {
    keyInfo: {
      type: Object,
      default: function() {
        return {
          serverId: 0,
          db: 0,
          key: ''
        }
      }
    }
  },
  data() {
    return {
      query: {
        serverId: 0,
        db: 0,
        type: '',
        key: '',
        cursor: 0,
        pageSize: 10,
        pageNo: 0
      },
      data: {
        type: '',
        key: '',
        ttl: -1,
        total: 0,
        value: []
      },
      loading: false
    }
  },
  mounted() {
    this.reload()
  },
  methods: {
    async reload() {
      this.loading = true
      const res = await keySummary(this.keyInfo)
      this.data = Object.assign({ key: this.keyInfo.type, value: [] }, res)
      this.query = Object.assign({ type: res.type }, this.keyInfo)
      switch (res.type.toUpperCase()) {
        case 'STRING':
          break
        case 'LIST':
          Object.assign(this.query, { pageNo: 1, pageSize: 10 })
          break
        case 'SET':
        case 'ZSET':
        case 'HASH':
          break
        default:
          this.loading = false
          return
      }
      this.loadNext()
      this.loading = false
    },
    async loadNext(param) {
      this.loading = true
      param = param ? Object.assign(Object.assign({}, this.query), param) : this.query
      const res = await getValue(param)
      switch (this.data.type.toUpperCase()) {
        case 'STRING':
          this.data.value = [{ value: res.value }]
          break
        case 'LIST':
          this.data.value = res.value.map(v => { return { value: v } })
          this.data.total = res.total
          this.query.pageSize = res.pageSize
          this.query.pageNo = res.pageNo
          break
        case 'SET':
          this.data.value = res.value.map(v => { return { value: v } })
          this.data.total = res.total
          this.query.cursor = res.cursor
          this.data.value = res.value.map(a => { return { value: a } })
          break
        case 'ZSET':
          this.data.value = []
          for (let i = 0; i < res.value.length;) {
            this.data.value.push({
              score: res.value[i++],
              value: res.value[i++]
            })
          }
          this.data.total = res.total
          this.query.cursor = res.cursor
          break
        case 'HASH':
          this.data.value = []
          for (let i = 0; i < res.value.length;) {
            this.data.value.push({
              key: res.value[i++],
              value: res.value[i++]
            })
          }
          this.data.total = res.total
          this.query.cursor = res.cursor
          break
        default:
          break
      }
      this.loading = false
    },
    async loadListPrev() {
      const param = Object.assign({}, this.query)
      param.pageNo--
      this.loadNext(param)
    },
    loadListNext() {
      const param = Object.assign({}, this.query)
      param.pageNo++
      this.loadNext(param)
    },
    openValueEditDialog(row) {
      this.$refs['valueEditDialog'].open(row)
    }
  }
}
</script>

<style scoped>
</style>
