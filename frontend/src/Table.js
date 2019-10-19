import React, { Component } from 'react'

const TableHeader = () => {
    return (
      <thead>
        <tr>
          <th>Name</th>
          <th>Username</th>
          <th>Bio</th>
          <th>Image</th>
        </tr>
      </thead>
    )
  }
    const TableBody = props => {
        const rows = props.characterData.map((row, index) => {
          return (
            <tr key={index}>
              <td>{row.name}</td>
              <td>{row.username}</td>
              <td>{row.bio}</td>
              <td>{row.authorimage}"</td>
            </tr>
          )
        })
    return (
      <tbody>
          {rows}
      </tbody>
    )
  }
  class Table extends Component {
    render() {
      const { characterData } = this.props
  
      return (
        <table>
          <TableHeader />
          <TableBody characterData={characterData} />
        </table>
      )
    }
  }

export default Table