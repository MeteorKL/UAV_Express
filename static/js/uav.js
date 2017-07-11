// http://blog.fooleap.org/bmaps-lushu.html

const p1 = new BMap.Point(120.131658, 30.273568) // 二食堂
const p2 = new BMap.Point(120.132193, 30.275040) // 30舍
const p3 = new BMap.Point(120.131402, 30.275439) // 31舍
var map
var uavs = []
var uavIcon
var uavPolylines = []
var uavMarks = []

function initMap(container) {
    map = new BMap.Map("container")          // 创建地图实例  
    map.addControl(new BMap.NavigationControl())
    map.addControl(new BMap.NavigationControl()) 
    map.addControl(new BMap.ScaleControl()) 
    map.addControl(new BMap.OverviewMapControl())
    map.centerAndZoom(p1, 18)
    uavIcon = new BMap.Icon('/img/uav.png', new BMap.Size(64, 40), {anchor: new BMap.Size(32, 20)})//动车
}

// function updateUav(uav) {
//     // console.log('updateUav', uav.ID)
//     pos = new BMap.Point(uav.Longitude, uav.Latitude)
//     uavMark[uav.ID].setPosition(pos);
// }

// function removeUav(uav) {
//     map.removeOverlay(uavMark[uav.ID]);
//     map.removeOverlay(uavPolyline[uav.ID]);
// }

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

function removeUav(id) {
    if (uavs[id] != null) {
        console.log("remove")
        map.removeOverlay(uavMarks[id])
        map.removeOverlay(uavPolylines[id])
    }
}

function updateUav(id, longitude, latitude, from_longitude, from_latitude, to_longitude, to_latitude) {
    if (uavs[id] == null) {
        console.log("new")
        from = new BMap.Point(from_longitude, from_latitude)
        to = new BMap.Point(to_longitude, to_latitude)
        var points = [from, to]
        uavPolylines[id] = new BMap.Polyline(points)
        map.addOverlay(uavPolylines[id])
        pos = new BMap.Point(longitude, latitude)
        uavMarks[id] = new BMap.Marker(pos)
        uavMarks[id].setPosition(pos)
        map.addOverlay(uavMarks[id])
    } else if (uavs[id].from_longitude!=from_longitude && uavs[id].from_latitude && uavs[id].to_longitude && uavs[id].to_latitude) {
        console.log("refresh")
        removeUav(id)
        updateUav(id, longitude, latitude, from_longitude, from_latitude, to_longitude, to_latitude)
    } else if (uavs[id].longitude!=longitude && uavs[id].latitude!=latitude) {
        console.log("move")
        pos = new BMap.Point(longitude, latitude)
        uavMarks[id].setPosition(pos)
    }
}
