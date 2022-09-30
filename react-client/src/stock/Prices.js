import React,{useState} from "react";
import JsonData from '../data/codes.json'
import { useTable, usePagination,useFilters } from 'react-table'
import {useNavigate} from "react-router-dom";
import 'bootstrap/dist/css/bootstrap.min.css';
import "./Prices.css"

function Table({ columns, data }) {
    const [filterInput, setFilterInput] = useState("");
    const navigate = useNavigate();
    const onRowClick = (e) => {
        navigate(`/${e.symbol}`);
    };

    const handleFilterChange = (e) => {
      const value = e.target.value || undefined;
      setFilter("symbol", value); 
      setFilterInput(value);
    };
    

    const {
        getTableProps,
        getTableBodyProps,
        headerGroups,
        prepareRow,
        setFilter,
        page,
        canPreviousPage,
        canNextPage,
        pageOptions,
        pageCount,
        gotoPage,
        nextPage,
        previousPage,
        setPageSize,
        state: { pageIndex, pageSize },
    } = useTable(
        {
            columns,
            data,
            initialState: { pageIndex: 0, pageSize: 5 },
        },
        useFilters,
        usePagination
    )

    // Render the UI for your table
    return (
        <div>
            <h1><span className="blue"></span>Current Stock<span className="blue"></span> <span className="yellow">Prices</span></h1>
            <h2>Created by <a href="https://github.com/VrMolodyakov" target="_blank">Vyachesav</a></h2>
            <input className = "symbol-input" value={filterInput} onChange={handleFilterChange} placeholder={"Search symbol"}/>
            <table className="container" {...getTableProps()}>
                <thead>
                    {headerGroups.map(headerGroup => (
                        <tr {...headerGroup.getHeaderGroupProps()}>
                            {headerGroup.headers.map(column => (
                                <th {...column.getHeaderProps()}>{column.render('Header')}</th>
                            ))}
                        </tr>
                    ))}
                </thead>
                <tbody {...getTableBodyProps()}>
                    {page.map((row, i) => {
                        prepareRow(row)
                        return (
                            <tr {...row.getRowProps()} onClick={() => onRowClick(row.original)}>
                                {row.cells.map(cell => {
                                    return <td {...cell.getCellProps()}>{cell.render('Cell')}</td>
                                })}
                            </tr>
                        )
                    })}
                </tbody>
            </table>
            {
            }
            <ul className="pagination">
                <li className="page-item" onClick={() => gotoPage(0)} disabled={!canPreviousPage}>
                    <a className="page-link">First</a>
                </li>
                <li className="page-item" onClick={() => previousPage()} disabled={!canPreviousPage}>
                    <a className="page-link">{'<'}</a>
                </li>
                <li className="page-item" onClick={() => nextPage()} disabled={!canNextPage}>
                    <a className="page-link">{'>'}</a>
                </li>
                <li className="page-item" onClick={() => gotoPage(pageCount - 1)} disabled={!canNextPage}>
                    <a className="page-link">Last</a>
                </li>
                <li>
                    <a className="page-link">
                        Page{' '}
                        <strong>
                            {pageIndex + 1 } of {pageOptions.length}
                        </strong>{' '}
                    </a>
                </li>
                <li>
                    <a className="page-link">
                        <input
                            className="form-control-inpt"
                            type="number"
                            defaultValue={pageIndex + 1 }
                            onChange={e => {
                                const page = e.target.value ? Number(e.target.value) - 1 : 0
                                gotoPage(page)
                            }}
                            style={{ width: '100px', height: '20px' }}
                        />
                    </a>
                </li>{' '}
                <select
                    className="form-control-inpt"
                    value={pageSize}
                    onChange={e => {
                        setPageSize(Number(e.target.value))
                    }}
                    style={{ width: '120px', height: '38px' }}
                >
                    {[5, 10, 20, 30, 40, 50].map(pageSize => (
                        <option key={pageSize} value={pageSize}>
                            Show {pageSize}
                        </option>
                    ))}
                </select>
            </ul>
        </div >
    )
}

function Prices() {
    const columns = React.useMemo(
        () => [
            {
                Header: ' ',
                columns: [
                    {
                        Header: "Symbol",
                        accessor: "symbol"
                    },
                    {
                        Header: "Name",
                        accessor: "name"
                    },
                    {
                        Header: "Exchange",
                        accessor: "exchange"
                    }
                ]
            }
        ],
        []
    )
    return (
        <Table columns={columns} data={JsonData} />
    )
}

export default Prices;