import React from 'react';
import axios from 'axios';

export default class Container extends React.Component {
    constructor(props) {
        super(props)
        this.state = {
            enterTime:0,
            appendHover:false,
            moving:-1,
            movingHover:-1,
            movingHoverBelow:false,
            hover:-1,
            focus:-1,
            focusChange:false,
            items:[
            ],
            movingItems:[]
        }
    }

    componentDidMount() {
        axios.get("http://ec2-34-210-40-109.us-west-2.compute.amazonaws.com/items/query")
        .then(res=> {
            console.log(res)
            this.setState({items:res.data.items});
        })
    }

    handleEnter(index,e) {
        var d = new Date();
        var n = d.getTime();
        console.log("Enter:",index)
        this.setState(prevState => ({
            hover:index,
            enterTime:n,
        }));
    }
    handleExit() {
        console.log("Exit")
        this.setState(prevState => ({
            hover:-1,
        }))
    }
    handleTextClick(index,e) {
        console.log("TextClick")
        this.setState(prevState => ({
            focus:index,
        }))
    }
    handleTextBlur() {
        console.log("TextBlur")
        if (this.state.focusChange) {
            axios.post("http://ec2-34-210-40-109.us-west-2.compute.amazonaws.com/items/update",{
                items:this.state.items,
            })
        }
        this.setState(prevState => ({
            focus:-1,
            focusChange:false
        }))
    }
    handleTextChange(index,e) {
        console.log("TextChange:",e.target)
        let newItems = this.state.items
        newItems[index]=e.target.value
        this.setState(prevState =>({
            items:newItems,
            focusChange:true
        }));
    }
    handleXClick(index,e) {
        console.log("xClick")
        var d = new Date();
        var n = d.getTime();
        if (n - this.state.enterTime  > 50 && index === this.state.hover) {
            let newItems = this.state.items
            newItems.splice(index,1)
            this.setState(prevState => ({
                items:newItems,
                hover:-1,
                focus:-1,
            }))
            axios.post("http://ec2-34-210-40-109.us-west-2.compute.amazonaws.com/items/update",{
                items:this.state.items,
            })
        }
    }
    handleAppendEnter(){
         console.log("AppendEnter")
        this.setState(prevState => ({
            appendHover:true,
            
        }))
    }
    handleAppendLeave(){
        console.log("AppendLeave")
        this.setState(prevState => ({
            appendHover:false,
        }))
    }
    handleAppendClick(){
        console.log("AppendClick")
        axios.post("http://ec2-34-210-40-109.us-west-2.compute.amazonaws.com/items/update",{
            items:this.state.items,
        })
        let newItems = this.state.items
        newItems.push("")
        this.setState(prevState =>({
            items:newItems,
        }))
    }
    handleMovingClick(index,e){
        console.log("MovingClick")
        var d = new Date();
        var n = d.getTime();
        if (n - this.state.enterTime  > 50 && index === this.state.hover) {
            this.setState(prevState => ({
                moving:index,
            }))
        }
    }
    handleMovingEnter(index,e){
        console.log("MovingEnter:A:",this.state.moving)
        let index1 = this.state.moving
        let index2 = index

        let newItems = this.state.items
        let temp = newItems.splice(index1,1)[0]
        console.log("MovingEnter:B:",temp)
        newItems.splice(index2,0,temp)
        this.setState(prevState =>({
            moving:index2,
            items:newItems,
        }))
    }
    handleMovingFinish(){
        console.log("MovingFinish")
        axios.post("http://ec2-34-210-40-109.us-west-2.compute.amazonaws.com/items/update",{
            items:this.state.items,
        })
        this.setState(prevState =>({
            hover:this.state.moving,
            moving:-1,
        }))
    }


    render () {
        console.log("Render Container")
        const items = this.state.items.map((item,index) =>
        <Box 
            key={index}
            value={item}
            handleEnter={this.handleEnter.bind(this,index)}
            handleExit={this.handleExit.bind(this)}
            handleTextClick={this.handleTextClick.bind(this,index)}
            handleTextBlur={this.handleTextBlur.bind(this)}
            handleTextChange={this.handleTextChange.bind(this,index)}
            handleXClick={this.handleXClick.bind(this,index)}
            handleMovingClick={this.handleMovingClick.bind(this,index)}
            handleMovingEnter={this.handleMovingEnter.bind(this,index)}
            handleMovingFinish={this.handleMovingFinish.bind(this)}

            hover={this.state.hover===index||this.state.focus===index}
            focus={this.state.focus===index}
            hoverIndex={this.state.hover}
            index={index}
            moving={this.state.moving===index}
            movingIndex={this.state.moving}
        />
        )
        let append = <div/>
        if (this.state.moving===-1) {
            append = <Append 
                hover={this.state.appendHover}
                handleAppendEnter={this.handleAppendEnter.bind(this)}
                handleAppendLeave={this.handleAppendLeave.bind(this)}
                handleAppendClick={this.handleAppendClick.bind(this)}
            />
        }
        return <div
            className="container"
        >
            <div
                className="box header"
            >Player Names</div>
            {items}
            {append}
        </div>
    }
}

function Append(props) {
    let className = "append"
    if (props.hover) {
        className = "appendHover"
    }
    return <div
        className={className}
        onMouseEnter={props.handleAppendEnter}
        onMouseLeave={props.handleAppendLeave}
        onClick={props.handleAppendClick}
    >
        <div 
            className="center"
        >
            <div
                className="octicon octicon-plus noselect "
                tabIndex="-1"
            />
        </div>
    </div>
}

function Box(props) {
    if (props.moving) {
        return <div>
                <div
                    className="box movingPreview"
                    onClick={props.handleMovingFinish}
                >{props.value}</div>
            </div> 
    } else if (props.movingIndex > -1) {
        return (
            <div
                className='box greyText'
                onMouseEnter={props.handleMovingEnter}
                onClick={props.handleMovingEnter}
            >{props.value}</div>
        )
    } else if (props.hover) {
        return <HoverBox
            handleExit={props.handleExit}
            handleTextClick={props.handleTextClick}
            handleTextBlur={props.handleTextBlur}
            handleXClick={props.handleXClick}
            handleTextChange={props.handleTextChange}
            handleMovingClick={props.handleMovingClick}
            value={props.value}
            focus={props.focus}
        />
    } else {
        return (
            <div
                className='box greyText'
                onMouseEnter={props.handleEnter}
                onClick={props.handleEnter}
            >{props.value}</div>
        )
    }
}

function HoverBox(props) {
    let textClass = "text"

    let textValue = props.value
    if (textValue==="" && !props.focus) {
        textValue = "Enter Value"
        textClass ="text lightGreyText"
    } 
    return (
        <div
            className='box hoverBox'
            onMouseLeave={props.handleExit}
        >
            <input 
                className={textClass}
                type="text"
                value={textValue}
                onClick={props.handleTextClick}
                onBlur={props.handleTextBlur}
                onChange={props.handleTextChange}
            />
            <div
                className="button xButton"
                onClick={props.handleXClick}
            >
                <div
                    className="octicon octicon-x noselect x"
                    tabIndex="-1"
                />
            </div>
            <div
                className="button orderButton"
                onClick={props.handleMovingClick}
            >
                <div 
                    className="octicon octicon-list-ordered noselect order"
                    tabIndex="-1"
                />
            </div>
        </div>
    )
}