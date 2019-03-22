using UnityEngine;
using System.Collections;
using System.Collections.Generic;
using UnityEngine.UI;


public class RoomPanel : PanelBase
{
    private List<Transform> prefabs = new List<Transform>();
    private Button closeBtn;
    private Button startOrPreBtn;                           //开始游戏或者准备按钮
    private Button joinRedTeamBtn;                          //加入红队按钮
    private Button joinBlueTeamBtn;                         //加入蓝队按钮
    private int isPrepare;                                  //玩家是否准备状态 0未准备 1准备
    public enum Team{ RED= 1,BLUE=2};                       //玩家队伍编号
    #region 生命周期
    /// <summary> 初始化 </summary>
    public override void Init(params object[] args)
    {
        base.Init(args);
        skinPath = "RoomPanel";
        layer = PanelLayer.Panel;
        isPrepare = 0;
    }

    public override void OnShowing()
    {
        base.OnShowing();
        Transform skinTrans = skin.transform;
        //组件
        for (int i = 0; i < 6; i++)
        {
            string name = "PlayerPrefab" + i.ToString();
            Transform prefab = skinTrans.Find(name);
            prefabs.Add(prefab);
        }
        closeBtn = skinTrans.Find("CloseBtn").GetComponent<Button>(); 
        startOrPreBtn = skinTrans.Find("StartOrPreBtn").GetComponent<Button>();
        joinRedTeamBtn = skinTrans.Find("JoinRedTeamBtn").GetComponent<Button>();
        joinBlueTeamBtn = skinTrans.Find("JoinBlueTeamBtn").GetComponent<Button>();
        //按钮事件
        closeBtn.onClick.AddListener(OnCloseClick);
        joinRedTeamBtn.onClick.AddListener(() => { OnJoinTeamClick(Team.RED); });
        joinBlueTeamBtn.onClick.AddListener(() => { OnJoinTeamClick(Team.BLUE); });
        //监听
        NetMgr.srvConn.msgDist.AddListener("GetRoomInfo", RecvGetRoomInfo);
        NetMgr.srvConn.msgDist.AddListener("Fight", RecvFight);
        //发送查询
        ProtocolBytes protocol = new ProtocolBytes();
        protocol.AddString("GetRoomInfo");
        NetMgr.srvConn.Send(protocol);
    }

    public override void OnClosing()
    {
        NetMgr.srvConn.msgDist.DelListener("GetRoomInfo", RecvGetRoomInfo);
        NetMgr.srvConn.msgDist.DelListener("Fight", RecvFight);
    }

    public void RecvGetRoomInfo(ProtocolBase protocol)
    {
        //获取总数
        ProtocolBytes proto = (ProtocolBytes)protocol;
        int start = 0;
        string protoName = proto.GetString(start, ref start);
        int count = proto.GetInt(start, ref start);
        //每个处理
        int i = 0;
        bool isMyselfOwner = false;
        for (i = 0; i < count; i++)
        {
            string id = proto.GetString(start, ref start);
            int team = proto.GetInt(start, ref start);
            int win = proto.GetInt(start, ref start);
            int fail = proto.GetInt(start, ref start);
            int isOwner = proto.GetInt(start, ref start);
            int isPrepare = proto.GetInt(start, ref start);
            //信息处理
            Transform trans = prefabs[i];
            Text text = trans.Find("Text").GetComponent<Text>();
            string str = "名字：" + id + "\r\n";
            str += "阵营：" + (team == 1 ? "红" : "蓝") + "\r\n";
            str += "胜利：" + win.ToString() + "   ";
            str += "失败：" + fail.ToString() + "\r\n";
            if (id == GameMgr.instance.id)
                str += "【我自己】";

            if (isOwner == 1)
            {
                str += "【房主】";
                if (id == GameMgr.instance.id)
                    isMyselfOwner = true;
            }
            else if (isPrepare == 1)
                str += " 已准备";
            else
                str += " 未准备";
            
            text.text = str;

            if (team == 1)
                trans.GetComponent<Image>().color = Color.red;
            else
                trans.GetComponent<Image>().color = Color.blue;
        }

        for (; i < 6; i++)
        {
            Transform trans = prefabs[i];
            Text text = trans.Find("Text").GetComponent<Text>();
            text.text = "【等待玩家】";
            trans.GetComponent<Image>().color = Color.gray;
        }
        //根据是否是房主 动态改变按钮的监听方法和按钮文字
        Text startText = startOrPreBtn.transform.Find("Text").GetComponent<Text>();
        startOrPreBtn.onClick.RemoveAllListeners();
        if (isMyselfOwner)
        {
            startOrPreBtn.onClick.AddListener(OnStartClick);
            startText.text = "开始战斗";
        }
        else
        {
            startOrPreBtn.onClick.AddListener(OnPrepareClick);
            startText.text = "准备";
        }
    }

    public void OnCloseClick()
    {
        ProtocolBytes protocol = new ProtocolBytes();
        protocol.AddString("LeaveRoom");
        NetMgr.srvConn.Send(protocol, OnCloseBack);
    }

    public void OnCloseBack(ProtocolBase protocol)
    {
        //获取数值
        ProtocolBytes proto = (ProtocolBytes)protocol;
        int start = 0;
        string protoName = proto.GetString(start, ref start);
        int ret = proto.GetInt(start, ref start);
        //处理
        if (ret == 0)
        {
            PanelMgr.instance.OpenPanel<TipPanel>("", "退出成功!");
            PanelMgr.instance.OpenPanel<RoomListPanel>("");
            Close();
        }
        else
        {
            PanelMgr.instance.OpenPanel<TipPanel>("", "退出失败！");
        }
    }

    public void OnStartClick()
    {
        ProtocolBytes protocol = new ProtocolBytes();
        protocol.AddString("StartFight");
        NetMgr.srvConn.Send(protocol, OnStartBack);
    }

    public void OnStartBack(ProtocolBase protocol)
    {
        //获取数值
        ProtocolBytes proto = (ProtocolBytes)protocol;
        int start = 0;
        string protoName = proto.GetString(start, ref start);
        int ret = proto.GetInt(start, ref start);
        //处理
        if (ret != 0)
        {
            PanelMgr.instance.OpenPanel<TipPanel>("", "开始游戏失败！两队至少都需要一名玩家，所有人都需要准备！");
        }
    }
    //准备按钮方法
    public void OnPrepareClick()
    {
        ProtocolBytes protocol = new ProtocolBytes();
        //根据是否准备，发送不同的消息
        if(isPrepare == 1)
            protocol.AddString("Cancel");
        else
            protocol.AddString("Prepare");
        NetMgr.srvConn.Send(protocol, OnPrepareBack);
    }

    public void OnPrepareBack(ProtocolBase protocol)
    {
        //获取数值
        ProtocolBytes proto = (ProtocolBytes)protocol;
        int start = 0;
        string protoName = proto.GetString(start, ref start);
        isPrepare = proto.GetInt(start, ref start);
        //根据是否准备 设置不同的按钮文字
        Text text = startOrPreBtn.transform.Find("Text").GetComponent<Text>();
        if (isPrepare == 1)
            text.text = "取消准备";
        else
            text.text = "准备";
    }
    //加入新的队伍方法
    public void OnJoinTeamClick(Team team)
    {
        //要求取消准备才能加入新的队伍
        if(isPrepare == 1)
        {
            PanelMgr.instance.OpenPanel<TipPanel>("", "请取消准备后更换队伍");
            return;
        }

        ProtocolBytes protocol = new ProtocolBytes();
        protocol.AddString("SwitchTeam");
        protocol.AddInt((int)team);
        NetMgr.srvConn.Send(protocol, OnJoinTeamBack);
    }

    public void OnJoinTeamBack(ProtocolBase protocol)
    {
    }

    public void RecvFight(ProtocolBase protocol)
    {
        ProtocolBytes proto = (ProtocolBytes)protocol;
        MultiBattle.instance.StartBattle(proto);
        PanelMgr.instance.OpenPanel<TalkPanel>("");
        Close();
    }

    #endregion
}