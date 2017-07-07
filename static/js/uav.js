// http://blog.fooleap.org/bmaps-lushu.html

const p1 = new BMap.Point(120.131658, 30.273568) // 二食堂
const p2 = new BMap.Point(120.132193, 30.275040) // 30舍
const p3 = new BMap.Point(120.131402, 30.275439) // 31舍
var map
var uavs = {}
var uavIcon
var uavPolyline = {}
var uavMark = {}
function initMap(container) {
    map = new BMap.Map("container")          // 创建地图实例  
    map.addControl(new BMap.NavigationControl())
    map.addControl(new BMap.NavigationControl()) 
    map.addControl(new BMap.ScaleControl()) 
    map.addControl(new BMap.OverviewMapControl())
    map.centerAndZoom(p1, 18)
    uavIcon = new BMap.Icon('/img/uav.png', new BMap.Size(64, 40), {anchor: new BMap.Size(32, 20)})//动车
}

function updateUav(uav) {
    // console.log('updateUav', uav.ID)
    pos = new BMap.Point(uav.Longitude, uav.Latitude)
    uavMark[uav.ID].setPosition(pos);
}

function removeUav(uav) {
    map.removeOverlay(uavMark[uav.ID]);
    map.removeOverlay(uavPolyline[uav.ID]);
}

function newUav(uav) {
    // console.log('newUav', uav.ID)
    from = new BMap.Point(uav.From.Longitude, uav.From.Latitude)
    to = new BMap.Point(uav.To.Longitude, uav.To.Latitude)
    pos = new  BMap.Point(uav.Longitude, uav.Latitude)
	var driving2 = new BMap.DrivingRoute(map, {renderOptions:{map: map, autoViewport: true}});    //驾车实例
    var points = [from, to]
    uavPolyline[uav.ID] = new BMap.Polyline(points)
    map.addOverlay(uavPolyline[uav.ID]);
    uavMark[uav.ID] = new BMap.Marker(pos)
    uavMark[uav.ID].setPosition(pos);
    map.addOverlay(uavMark[uav.ID]);
}
