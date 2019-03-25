using UnityEngine;
using System.Collections;
using UnityEngine.UI;

public class LoginPanel : PanelBase
{
    private InputField idInput;
    private InputField pwInput;
    private Button loginBtn;
    private Button regBtn;
    private string newHandMessage;

    #region 生命周期
    //初始化
    public override void Init(params object[] args)
    {
        base.Init(args);
        skinPath = "LoginPanel";
        layer = PanelLayer.Panel;
        newHandMessage = "您好，尊敬的坦克兵！在踏上保卫祖国的旅途前，你需要了解一下操作方法：\n " + 
            "用鼠标控制坦克的射击，W，A，S，D控制坦克的移动，回车键进行聊天的选择与发送.祝你好运！";
    }

    public override void OnShowing()
    {
        base.OnShowing();
        Transform skinTrans = skin.transform;
        idInput = skinTrans.Find("IDInput").GetComponent<InputField>();
        pwInput = skinTrans.Find("PWInput").GetComponent<InputField>();
        loginBtn = skinTrans.Find("LoginBtn").GetComponent<Button>();
        regBtn = skinTrans.Find("RegBtn").GetComponent<Button>();

        loginBtn.onClick.AddListener(OnLoginClick);
        regBtn.onClick.AddListener(OnRegClick);
    }
    #endregion

    public void OnRegClick()
    {
        PanelMgr.instance.OpenPanel<RegPanel>("");
        //Close();
    }

    public void OnLoginClick()
    {
        //用户名密码为空
        if (idInput.text == "" || pwInput.text == "")
        {
            PanelMgr.instance.OpenPanel<TipPanel>("", "用户名密码不能为空!");
            return;
        }
        //连接服务器
        if (NetMgr.srvConn.status != Connection.Status.Connected)
        {
            string host = "10.246.34.154"; 
            int port = 18085; 
            NetMgr.srvConn.proto = new ProtocolBytes();
            if (!NetMgr.srvConn.Connect(host, port))
                PanelMgr.instance.OpenPanel<TipPanel>("", "连接服务器主端口失败!");
        }
        //发送
        ProtocolBytes protocol = new ProtocolBytes();
        protocol.AddString("Login");
        protocol.AddString(idInput.text);
        protocol.AddString(pwInput.text);
        protocol.AddString("xinjaystudio");
        Debug.Log("发送 " + protocol.GetDesc());
        NetMgr.srvConn.Send(protocol, OnLoginBack);


        //连接服务器
        if (NetMgr.talkConn.status != Connection.Status.Connected)
        {
            string host = "10.246.34.154";
            int port = 18086;
            NetMgr.talkConn.proto = new ProtocolBytes();
            if (!NetMgr.talkConn.Connect(host, port))
                PanelMgr.instance.OpenPanel<TipPanel>("", "连接服务器聊天端口失败!");
        }
        ProtocolBytes talkProtocol = new ProtocolBytes();
        talkProtocol.AddString(idInput.text);
        NetMgr.talkConn.Send(talkProtocol, OnTalkBack);
    }

    public void OnTalkBack(ProtocolBase protocol)
    {

    }

    public void OnLoginBack(ProtocolBase protocol)
    {
        ProtocolBytes proto = (ProtocolBytes)protocol;
        int start = 0;
        string protoName = proto.GetString(start, ref start);
        int ret = proto.GetInt(start, ref start);
        int isNewerHand = proto.GetInt(start, ref start);
        if (ret == 0)
        {
            if (isNewerHand == 1)
                PanelMgr.instance.OpenPanel<TipPanel>("", newHandMessage);
            else
                PanelMgr.instance.OpenPanel<TipPanel>("", "登录成功!");

            //开始游戏
            PanelMgr.instance.OpenPanel<RoomListPanel>("");
            GameMgr.instance.id = idInput.text;
            Close();
        }
        else
        {
            PanelMgr.instance.OpenPanel<TipPanel>("", "登录失败，请检查用户名密码或者口令!");
        }
    }
}