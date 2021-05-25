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

// type nodeData = {
//     key: string,
//     workerUrl: string,
//     id: number,
//     status: string,
//     amount: number
// }[]

const columns: missionNode = [

    {
        title: '节点编号',
        dataIndex: 'id',
        key: 'id',
    },
    {
        title: '节点IP',
        dataIndex: 'workerUrl',
        key: 'workerUrl',
    },
    {
        title: '当前状态',
        dataIndex: 'status',
        key: 'status',
    },
    {
        title: '当前任务数',
        dataIndex: 'amount',
        key: 'amount',
    },


];
// const data: nodeData = [
//     {
//         key: '1',
//         workerUrl: 'John Brown',
//         id: 32,
//         status: 'New York No. 1 Lake Park',
//         amount: 1,
//     },
//     {
//         key: '2',
//         workerUrl: 'Jim Green',
//         id: 42,
//         status: 'London No. 1 Lake Park',
//         amount: 2,
//     },
//     {
//         key: '3',
//         workerUrl: 'Joe Black',
//         id: 32,
//         status: 'Sidney No. 1 Lake Park',
//         amount: 1,
//     },

// ];


/** 组件 */
const Node: React.FC = () => {
    /** 分页 */
    const [pageSize, setPagesize] = useState(5)
    const [pageNum, setPagenum] = useState(1)
    const [total, setTotal] = useState(8)
    const [mydata, setMyData] = useState([])
    useEffect(()=>{
        (async function getWorkListData(){
            const response = await taskService.getWorkList(total, 0)
            console.log(response);
            
            setMyData(response.workerInfo)
            setTotal(response.count)
        })()
        
        
    }, [])
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
            <Table columns={columns} dataSource={mydata} pagination={paginationProps} />
        </div>
    )
}

export default Node
