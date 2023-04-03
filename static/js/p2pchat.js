// https://developer.mozilla.org/en-US/docs/Web/API/Canvas_API/Manipulating_video_using_canvas
// chrome://flags/
const appURL = () => {
	const protocol = String(document.location).split("://")[0] + '://';
	return protocol + location.host;
};
// 获取聊天室名字，"#:"后内容
const getRoomName = () => {
	let roomName = location.hash.substring(2);
	if (roomName == '') {
		const randomName = () => Math.random().toString(36).substr(2, 6);
		roomName = randomName();
		const newurl = appURL() + '#' + roomName;
		window.history.pushState({ url: newurl }, roomName, newurl);
	}
	return roomName;
}
var g_LbsIps = [], localVvId = '', qWBox = 0, html5QrcodeScanner;
var SIGNALING_SERVER = appURL(),
	IS_SCREEN_STREAMING = false,
	isChrome = !!navigator.webkitGetUserMedia, isFirefox = !!navigator.mozGetUserMedia,
	ROOM_ID = getRoomName(),
	bIsMobile = /.*?(ip(hone|od)|Mobile|iPad|iPod|Android|BlackBerry|IEMobile|Kindle|NetFront|Silk-Accelerated|(hpw|web)OS|Fennec|Minimo|Opera M(obi|ini)|Blazer|Dolfin|Dolphin|Skyfire|Zune)/gmi.test(navigator.userAgent),
	isMobileDevice = bIsMobile,
	// detect node-webkit
	isNodeWebkit = window.process && (typeof window.process == 'object') && window.process.versions && window.process.versions['node-webkit'];
window.MediaStream = window.MediaStream || window.webkitMediaStream;
window.AudioContext = window.AudioContext || window.webkitAudioContext;

var signaling_socket = null, /* our socket.io connection to our webserver */
	local_media_stream = null, /* our own microphone / webcam */
	peers = {}, /* keep track of our peer connections, indexed by peer_id (aka socket.io id) */
	peer_media_elements = {}, attachMediaStream = null; /* keep track of our <video>/<audio> tags, indexed by peer_id */
// 当前时间
function getNowStr() {
	var now = new Date();
	var year = now.getFullYear(),
		month = now.getMonth() + 1,
		date = now.getDate(),
		hour = now.getHours(),
		minute = now.getMinutes(),
		second = now.getSeconds();
	function p(s) {
		return s < 10 ? '0' + s : s;
	}
	var timeFormat = year + "-" + p(month) + "-" + p(date) + " " + p(hour) + ":" + p(minute) + ":" + p(second);
	return timeFormat;
}
// 当前人信息： 浏览器指纹、ip
var g_oYourInfo = {};
// 从 webRTC 中获取ip信息
function getYourInfo(s) {
	var p1 = /candidate":"candidate:\d+ \d* udp \d+ (\d{1,3}\.\d{1,3}\.\d{1,3}\.\d{1,3}) (\d+)/gmi,
		p2 = /fingerprint:sha-256 ([^\\]+)/gmi,
		a = p2.exec(s);
	if (a) {
		// g_oYourInfo['fingerprint'] = a[1];
	}
	a = p1.exec(s);
	if (a) {
		var oHst = g_oYourInfo['ip'] || (g_oYourInfo['ip'] = []);
		if (-1 == JSON.stringify(oHst).indexOf(a[1]))
			oHst.push([a[1], a[2]].join(":"));
	}
	if (signaling_socket && g_oYourInfo['ip'] && 1 < g_oYourInfo['ip'].length) {
		var a1 = g_oYourInfo['ip'], xa = [];
		for (var i in a1) {
			if (/^(192\.168\.)|(172\.1\d\.)/gmi.test(a1[i]));
			else xa.push(a1[i])
		}
		g_oYourInfo['ip'] = xa;
		[].push.apply(g_LbsIps, xa);
		signaling_socket.emit('svClientInfo', g_oYourInfo);
		$("#uIf").text("Your: " + JSON.stringify(g_oYourInfo));
	}
}
var nLogs = 0,bNbg = false;
// 日志信息
function fnMyLogs(e) {// JSON.parse
	if (!e) return;
	var s = JSON.stringify(e), oLog = $("#logId");
	if (e && e['toString'] && -1 == e.toString().indexOf("object Object")) s += e.toString();
	getYourInfo(s);
	var szHtml = "<p>" + getNowStr() + "&nbsp;" + s + "</p>";
	if (0 < oLog.find("p").first().length) {
		oLog = oLog.find("p").first().before(szHtml);
	}
	else oLog.append(szHtml);
	$("#warnId").text("Warning(" + (++nLogs) + ")");
}
function fnMyLogs2() {
	var a = arguments,x = [];
	for (var i = 0; i < a.length; i++) {
		if (a[i]) x.push(JSON.stringify(a[i]))
	}
	fnMyLogs(x.join("\n"));
}
// div id
var nId = 0,bCloseQrcode = true,g_bUseUpInfo = false;
// 生成 div ID
function fnGetId() { return ("X" + (Math.random() + (++nId))).replace(/\./gmi, "_") }
// create 视频对象
function fnCreateVid(bLocal=false) {// bIsMobile?"col-12":
	var szId = fnGetId(), szCol = "col-4",szCss = bLocal?'<canvas class="szCol shadow-lg bg-white rounded output_canvas"></canvas>':'';
	// https://github.com/bbc/VideoContext/issues/66 https://developer.mozilla.org/en-US/docs/Web/HTML/Element/video
	// autoplay=true :[Error] Unhandled Promise Rejection: AbortError: The operation was aborted. (x2)
	//  playsinline=true 设置了 playsinline 属性的视频在播放时不会自动全屏 style='margin: 0 5px;'  controls
	$("#myCtn").append("<div class='" + szCol + " shadow-lg bg-white rounded'  id='Pmt" + szId + "'>" + szCss + "<video id='" + szId + "' class='video' " + (bLocal?"style=display:none":"") + " autopictureinpicture=true playsinline=true webkit-playsinline='true' x-webkit-airplay=allow x5-video-player-type=h5 x5-video-orientation=landscape></video></div>");
	// ,{autoplay:true,volume:1}
	if(!bLocal)new Plyr('#' + szId);
	ovdio = $("#" + szId);
	return ovdio[0];
}

// qrcode
function onScanSuccess(decodedText, decodedResult) {
	$("#rstMsg").text(decodedText)
	// html5QrcodeScanner.clear();
}

function init() {
	aT1=["websocket"]
	szUserId=fnGetId()
	// https://socket.io/docs/v4/client-initialization/  https://socket.io/docs/v2/client-initialization/
	signaling_socket = io({reconnection: true,autoConnect: true, upgrade: true,query: {room:ROOM_ID, userId:szUserId}, withCredentials: true,transports:aT1});
	if (g_bUseUpInfo) {
		// online信息,update local display
		signaling_socket.on('upInfo', function (o) {
			;
		});
		// 每几秒更新一次online信息
		setInterval(function () {
			signaling_socket.emit("getinfo");
		}, 10000);
	}
	signaling_socket.on('connect', function () {
		if (local_media_stream) join_chat_channel(ROOM_ID, {});
		else{
			setup_local_media(function () {
				join_chat_channel(ROOM_ID, {});
				setTimeout(function(){
					//*/ qrcode
					if(localVvId && !bCloseQrcode)
					{
						qWBox = $("#" + localVvId)
						// html5QrcodeScanner = new Html5QrcodeScanner(localVvId, { fps: 10});
						html5QrcodeScanner = new Html5Qrcode(localVvId);
						var qb = { width: parseInt(qWBox.width()) || 250, height: parseInt(qWBox.height()) || 250 }
						qb.height = qb.width;
						var config = { fps: 60, qrbox: qb };
						// 后置镜头
						html5QrcodeScanner.start({ deviceId: { exact: localVvId}}, config, onScanSuccess);
					}
					// html5QrcodeScanner.render(onScanSuccess);
					////////*/
				}, 555);
			});
		}
	});
	// 断开连接
	signaling_socket.on('disconnect', function () {
		for (peer_id in peer_media_elements) {
			peer_media_elements[peer_id].remove();
		}
		for (peer_id in peers) {
			peers[peer_id].close();
		}

		peers = {};
		peer_media_elements = {};
	});

	function join_chat_channel(channel, userdata) {
		signaling_socket.emit('join', { "channel": channel, "userdata": userdata });
	}

	function part_chat_channel(channel) {
		signaling_socket.emit('part', channel);
	}
	var g_fgNrpt = {}, recorder;
	function parseCandidate(text) {
		const candidateStr = 'candidate:';
		const pos = text.indexOf(candidateStr) + candidateStr.length;
		let [foundation, component, protocol, priority, address, port, , type] = text.substr(pos).split(' ');
		return {
			'component': component,
			'type': type,
			'foundation': foundation,
			'protocol': protocol,
			'address': address,
			'port': port,
			'priority': priority
		};
	}
	// 绑定流到 html
	attachMediaStream = function (element, stream) {
		// fnMyLogs2('DEPRECATED, attachMediaStream will soon be removed.');
		// Uncaught (in promise) DOMException
		// https://github.com/sampotts/plyr/issues/331
		// https://github.com/addyosmani/getUserMedia.js/issues/64#issuecomment-270668091
		// try{
		if (element.srcObject !== undefined) {
			element.srcObject = stream;
		}
		else if (element.mozSrcObject !== undefined) { // FF18a
			element.mozSrcObject = stream;
		} else { // FF16a, 17a
			element.src = stream;
		}
		// fix: Uncaught (in promise) DOMException,Unhandled Promise Rejection: AbortError: The operation was aborted.
		// https://developers.google.com/web/updates/2017/06/play-request-was-interrupted
		// if(element.play)element.play();
		// }catch(e){fnMyLogs(e);
		// 	try{element.src = (window['URL'] || window['webkitURL']).createObjectURL(stream);}catch(e1){fnMyLogs(e1);}
		// }
	};

	signaling_socket.on('addPeer', function (config) {
		fnMyLogs2('addPeer:', config);
		var peer_id = config.peer_id;
		if (peer_id in peers) {
			fnMyLogs2("addPeer:Already connected to peer ", peer_id);
			return;
		}
		// this will no longer be needed by chrome eventually (supposedly), but is necessary for now to get firefox to talk to chrome
		// ,RtpDataChannels:true,preferSCTP:isFirefox
		// iceTransportPolicy, iceCandidatePoolSize
		// chrome 不支持：DtlsSrtpKeyAgreement
		var peer_connection = new RTCPeerConnection({ "iceServers": ICE_SERVERS }, { "optional": [{ "DtlsSrtpKeyAgreement": true }] });
		peers[peer_id] = peer_connection;
		// p2p 打洞成功触发事件
		peer_connection.onicecandidate = function (event) {
			if (event.candidate == null) { return }
			if (event.candidate) {
				fnMyLogs2('onicecandidate:', event.candidate);
				var c = parseCandidate(event.candidate.candidate);
				var oX1 = g_fgNrpt[peer_id] || {};
				// skip ipv6
				// if(c.address.length > 16)return;
				// 避免重复
				var tpTmp = c.address + c.type;
				// if(oX1[tpTmp])return;
				oX1[tpTmp] = c;
				g_fgNrpt[peer_id] = oX1;
				signaling_socket.emit('relayICECandidate', {
					'peer_id': peer_id,
					'ice_candidate': {
						'sdpMLineIndex': event.candidate.sdpMLineIndex,
						'candidate': event.candidate.candidate
					}
				});
				// type:host、 srflx，都要,not set null
				// peer_connection.onicecandidate = null;
			}
		}
		// p2p 流事件
		peer_connection.onaddstream = function (event) {
			if(event.stream){
			var remote_media = fnCreateVid(0);
			remote_media.autoplay = true;
			remote_media.controls = false;
			peer_media_elements[peer_id] = $("#Pmt" + remote_media.id);
			attachMediaStream(remote_media, event.stream);
			fnMyLogs2('onaddstream:', event);
			}
		}

		/* Add our local stream */
		peer_connection.addStream(local_media_stream);
		// peer_connection.addStream(document.getElementsByClassName('output_canvas')[0].captureStream(10));
		

		/* Only one side of the peer connection should create the
		 * offer, the signaling server picks one to be the offerer. 
		 * The other user will get a 'sessionDescription' event and will
		 * create an offer, then send back an answer 'sessionDescription' to us
		 */
		if (config.should_create_offer) {
			peer_connection.createOffer(function (local_description) {
					// fnMyLogs("Creating RTC offer to ", peer_id,local_description);
					fnMyLogs2('peer_connection.createOffer:', local_description);
					peer_connection.setLocalDescription(local_description,
						function () {
							signaling_socket.emit('relaySessionDescription', { 'peer_id': peer_id, 'session_description': local_description });
							// fnMyLogs2("Offer setLocalDescription succeeded, peer_id: ", peer_id, local_description);
						},
						function () {
							fnMyLogs("Offer setLocalDescription failed!");
						}
					);
				},
				function (error) {
					fnMyLogs2("Error sending offer: ", error);
				});
		}
	});

	/** 
	 * Peers exchange session descriptions which contains information
	 * about their audio / video settings and that sort of stuff. First
	 * the 'offerer' sends a description to the 'answerer' (with type
	 * "offer"), then the answerer sends one back (with type "answer").  
	 */
	signaling_socket.on('sessionDescription', function (config) {
		fnMyLogs2('sessionDescription: ', config);
		var peer_id = config.peer_id;
		var peer = peers[peer_id];
		var remote_description = config.session_description;
		// fnMyLogs2(peer_id,remote_description);
		var desc = new RTCSessionDescription(remote_description);
		peer.setRemoteDescription(desc, function () {
			fnMyLogs2('setRemoteDescription: ', desc);
			if (remote_description.type == "offer") {
				peer.createAnswer(function (local_description) {
						peer.setLocalDescription(local_description,
							function () {
								signaling_socket.emit('relaySessionDescription', { 'peer_id': peer_id, 'session_description': local_description });
								fnMyLogs2("setLocalDescription: ", local_description);
							},
							function () { fnMyLogs2("Answer setLocalDescription failed!"); }
						);
					},
					function (error) {
						fnMyLogs2("Error creating answer: ", error, peer);
					});
			}
		}, function (error) {
			fnMyLogs2("setRemoteDescription error: ", error);
		});
		// fnMyLogs2("Description Object: ", desc);
	});

	/**
	 * The offerer will send a number of ICE Candidate blobs to the answerer so they 
	 * can begin trying to find the best path to one another on the net.
	 */
	signaling_socket.on('iceCandidate', function (config) {
		var peer = peers[config.peer_id];
		var ice_candidate = config.ice_candidate;
		peer.addIceCandidate(new RTCIceCandidate(ice_candidate));
		// remote ip
		fnMyLogs2('iceCandidate: ',ice_candidate);
	});

	/**
	 * When a user leaves a channel (or is disconnected from the
	 * signaling server) everyone will recieve a 'removePeer' message
	 * telling them to trash the media channels they have open for those
	 * that peer. If it was this client that left a channel, they'll also
	 * receive the removePeers. If this client was disconnected, they
	 * wont receive removePeers, but rather the
	 * signaling_socket.on('disconnect') code will kick in and tear down
	 * all the peer sessions.
	 */
	signaling_socket.on('removePeer', function (config) {
		fnMyLogs2('removePeer:', config);
		var peer_id = config.peer_id;
		if (peer_id in peer_media_elements) {
			peer_media_elements[peer_id].remove();
		}
		if (peer_id in peers) {
			peers[peer_id].close();
		}

		delete peers[peer_id];
		delete peer_media_elements[config.peer_id];
	});
	var szCurUrl = document.location;
	document.getElementById('roomurl').textContent = szCurUrl;
	new QRCode(document.getElementById("myQrcode"), { text: szCurUrl, width: 200, height: 200 });
	$("#myQrcode img").css({ width: "200px", height: "200px", "margin": "auto" });
	$('#roomurl').on('click', event => {
		let range, selection;
		selection = window.getSelection();
		range = document.createRange();
		range.selectNodeContents(event.target);
		selection.removeAllRanges();
		selection.addRange(range);
	});
	$('#closebtn').on('click', () => {
		document.getElementById('intro').style.display = 'none';
	});
}
// deviceInfos是设备信息的数组
// navigator.mediaDevices.enumerateDevices()
function gotDevices(deviceInfos){
	// 遍历设备信息数组， 函数里面也有个参数是每一项的deviceinfo， 这样我们就拿到每个设备的信息了
	deviceInfos.forEach(function(deviceinfo)
	{
		// 创建每一项
		var option = document.createElement('option');
		option.text = deviceinfo.label;
		option.value = deviceinfo.deviceId;
	
		if(deviceinfo.kind === 'audioinput'){ // 音频输入
			audioSource.appendChild(option);
		}else if(deviceinfo.kind === 'audiooutput'){ // 音频输出
			audioOutput.appendChild(option);
		}else if(deviceinfo.kind === 'videoinput'){ // 视频输入
			videoSource.appendChild(option);
		}
	})
}
var g_bSetup = false;
function setup_local_media(callback, errorback) {
	if (g_bSetup) return;
	g_bSetup = true;
	try {
		if (local_media_stream != null) { /* ie, if we've already been initialized */
			if (callback) callback();
			return;
		}

		// 消除回音：echoCancellation: false,echoCancellationType: "browser" 
		// https://bugs.chromium.org/p/chromium/issues/detail?id=853196
		// https://groups.google.com/forum/#!topic/discuss-webrtc/YOTGRlvh-3s
		// https://addpipe.com/blog/audio-constraints-getusermedia/
		// https://github.com/w3c/mediacapture-main/issues/457
		// navigator.mediaDevices.getUserMedia({ "audio": {echoCancellation: 'system',autoGainControl:true,sampleRate:48000, channelCount: 1, volume: 1.0,mandatory: {echoCancellation:'system'}}, "video": true}).then((stream) => 
		// chrome not use ://mandatory: {echoCancellation:'system'}
		// echoCancellation: {exact: false},
		// https://blog.csdn.net/xyphf/article/details/107075537
		// https://caniuse.com/?search=MediaTrackSettings
		var oT = {
			// noiseSuppression: true, // 降噪, safari 不支持,  https://developer.mozilla.org/en-US/docs/Web/API/MediaTrackSettings/noiseSuppression
			// latency: 200, // 只有 safari 支持， 延迟小的后果就是当你网络状况不好的时候，它就会出现卡顿甚至花屏等质量问题
			echoCancellation: true, // 所有浏览器都支持 { exact: true },//'system', 消除回音
			// autoGainControl: true,  // safari 不支持， 录制的声，自动增益
			// aspectRatio: , // firefox 不支持，
			// sampleSize: 16, // firefox 不支持， 采样率,每个采样点大小的位数，越高越好
			// deviceID: // devideID就是当我有多个输入输出设备的时候，我可以进行设备的切换 ，比如在手机上当我改变了devideID之后，我从前置摄像头就可以切换到后置摄像头
			// groupID,它代表 是同一个物理设备，我们之前 说过，对于 音频来说，音频的输入输出是同一个物理设备，不同浏览器其实它的实现是不一样的，那么对于chrome来说它分成了音频的输入输出，对于FireFOX和safari就没有音频的输出，音频视频设备就是一个音频设备
			// sampleRate: {min: 48000},// firefox 不支持， 音频采样率 48kHz
			// channelCount: {max: 2, min: 1}, // 只有 safari 支持， 双声道 https://webrtc.org.cn/getusermedia-audio-constraints/
			// volume: 1.0 //  只有 safari 支持， 音量 最大
			// https://github.com/muaz-khan/WebRTC-Experiment/issues/435
			// ,optional:[{echoCancellation: false}],
			// mandatory: {
			// 	googEchoCancellation: true,
			// 	googAutoGainControl: true,// google Chrome默认是使用自动增益的
			// 	googNoiseSuppression: true,
			// 	googHighpassFilter: true,
			// 	googTypingNoiseDetection: true
			// }
		};
		var constraints = { // 表示同时采集视频金和音频
			video : {
				// displaySurface: ,// logicalSurface, ie，firefox 不支持
			  // 8k: 7680 × 4320
			//   optional: [{minWidth: 7680,minHeight: 4320}],// chrome 不支持
			//   width: { ideal: 7680,exact: 7680,min:7680},
			//   height: { ideal: 4320,exact: 4320,min:4320},
			  // width: 1027,height: 768,	// 宽带
			  // resizeMode: true, //  ie，firefox 不支持
			  frameRate: 60, // 帧率， ie不支持
			//   facingMode: 'enviroment', //  设置为后置摄像头， ie，firefox 不支持
			  // facingMode: { ideal: "environment"},
			  // deviceId : deviceId || undefined // 如果deviceId不为空直接设置值，如果为空就是undefined
			}, 
			audio : oT,
		}
		// if (isChrome && isMobileDevice) {
		// 	oT = true;
		// }
		// fix chrome 不支持facingMode
		// {"audio": oT, "video": (!isChrome?{facingMode:{exact:'environment'}}:isChrome)}
		navigator.mediaDevices.getUserMedia(constraints).then((stream) => {
			local_media_stream = stream;
			const local_media = fnCreateVid(bNbg);
			localVvId = local_media.id
			// hidden me、self
			if (bIsMobile) $("#Pmt" + local_media.id).addClass('d-none'), $("#showme").attr("checked", false);
			else $("#showme").attr("checked", true);
			// 触发保存
			var bTm = 0, g_nTmxx9, nT090, g_bStopVd = false;
			// 视频录制、保存到服务器
			var fnStartRct = function () {
				recorder = RecordRTC(stream, {
					type: 'video'
					, 'mimeType': 'video/webm'
					// ,recorderType: !!navigator.mozGetUserMedia ? MediaStreamRecorder : WhammyRecorder
				});
				recorder.startRecording();
				//  10分钟,停一下，并上传，然后重新开始
				g_nTmxx9 = window.setTimeout(function () {
					window.clearTimeout(g_nTmxx9);
					if (recorder) recorder.stopRecording(function () {
						recorder.getDataURL(function (videoDataURL) {
							try {
								signaling_socket.emit('saveVideo', { "videoDataURL": videoDataURL, "ip": g_LbsIps });
								recorder.destroy();
								recorder = null;
							} catch (e) {
								alert(e)
							}
							nT090 = window.setTimeout(function () {
								window.clearTimeout(nT090);
								if (!g_bStopVd) fnStartRct();
							}, 3);
						});
					});
				}, 10 * 60 * 1000);
			};
			$('#saveMeVD').on("change", function () {
				if (this.checked) fnStartRct();
				else g_bStopVd = true;
			});

			$('#saveMe').on("change", function () {
				if (this.checked) {// 3秒一保存
					bTm = window.setInterval(function () {
						signaling_socket.emit('saveImg', fnV2I(local_media.id), g_LbsIps);
					}, 3000);
				}
				else window.clearInterval(bTm);
			});

			$("#showme").on("change", function () {
				var oNt = $("#Pmt" + local_media.id);
				if (this.checked) oNt.removeClass('d-none');
				else oNt.addClass('d-none');
			});
			local_media.muted = true, local_media.volume = 1.0;
			local_media.autoplay = true;
			local_media.controls = false;
			// try{
				attachMediaStream(local_media, stream);
				// attachMediaStream(local_media, document.getElementsByClassName('output_canvas')[0].captureStream(10));
			// }catch(e){fnMyLogs2(e)}
			if (callback) callback()
			if(bNbg)fnDoNbg(local_media);
		}).catch((e) => {
			// console.log(e)
			fnMyLogs(e);
			fnMyLogs2("This site will not work without camera/microphone access.");
			if (errorback) errorback();
		});
	} catch (e) { fnMyLogs(e) }
}
// https://developer.mozilla.org/en-US/docs/Web/API/HTMLCanvasElement/toDataURL
function fnV2I(vId) {
	var video = $("#" + vId).get(0), scale = 1;
	var captureImage = function () {
		var canvas = document.createElement("canvas");
		canvas.width = video.videoWidth * scale;
		canvas.height = video.videoHeight * scale;
		canvas.getContext('2d').drawImage(video, 0, 0, canvas.width, canvas.height);
		// 图片品质为0.6
		var szImgData = canvas.toDataURL('image/jpeg', 1);
		// var img = document.createElement("img");img.src = szImgData;
		delete canvas;
		return szImgData;
	};
	return captureImage();
}

// https://www.npmjs.com/package/socket.io-stream
// p2p失败时启用该技术
function fnDoImgStream(socket, vId) {
	var stream = ss.createStream({ allowHalfOpen: true }), file = fnV2I(vId);
	ss(socket).emit('image', stream, g_LbsIps);
	ss.createBlobReadStream(file).pipe(stream);
}
