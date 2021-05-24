import React, { useEffect, useState } from 'react'
import { Table } from 'antd';
import './style.less'
import { task as taskService } from '@/services';


/** 定义节点类型 */
type missionNode = {
    title: string,
    dataIndex: string,
    key: string
}[]

type nodeData = {
    key: string,
    nodeIp: string,
    nodeOrder: number,
    currentState: string,
    amount: number
}[]

const columns: missionNode = [

    {
        title: '节点编号',
        dataIndex: 'nodeOrder',
        key: 'nodeOrder',
    },
    {
        title: '节点IP',
        dataIndex: 'nodeIp',
        key: 'nodeIp',
    },
    {
        title: '当前状态',
        dataIndex: 'currentState',
        key: 'currentState',
    },
    {
        title: '当前任务数',
        dataIndex: 'amount',
        key: 'amount',
    },


];
const data: nodeData = [
    {
        key: '1',
        nodeIp: 'John Brown',
        nodeOrder: 32,
        currentState: 'New York No. 1 Lake Park',
        amount: 1,
    },
    {
        key: '2',
        nodeIp: 'Jim Green',
        nodeOrder: 42,
        currentState: 'London No. 1 Lake Park',
        amount: 2,
    },
    {
        key: '3',
        nodeIp: 'Joe Black',
        nodeOrder: 32,
        currentState: 'Sidney No. 1 Lake Park',
        amount: 1,
    },

];


/** 组件 */
const Node: React.FC = () => {
    /** 分页 */
    const [pageSize, setPagesize] = useState(5)
    const [pageNum, setPagenum] = useState(1)
    const [total] = useState(8)
    useEffect(()=>{
        async function getWorkListData(){
            await taskService.getWorkList(total, pageSize)
        }
        getWorkListData()
        
    })
    const changePage = (n: any, s: any) => {
        setPagenum(n)
        setPagesize(s)
    }

    const paginationProps = {
        showSizeChanger: true,
        showQuickJumper: false,
        showTotal: () => `共${total}条`,
        pageSize,
        current: pageNum,
        total,
        onChange: changePage,
    }
    return (
        <div>
            <p className='nodeTitle'>任务执行节点列表</p>
            <Table columns={columns} dataSource={data} pagination={paginationProps} />
        </div>
    )
}

export default Node
