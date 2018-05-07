#ifndef __PAAS_AES_CRYPTO_H__
#define __PAAS_AES_CRYPTO_H__

#define PAAS_CRYPTO_PATH       "/var/paas"
#define PRIMARY_KEYSTORE_FILE  "/var/paas/primary_store.txt"
#define STANDBY_KEYSTORE_FILE  "/var/paas/standby_store.txt"

#define PRIMARY_KMC_CFG_FILE  "/var/paas/primary_kmc_cfg.txt"
#define STANDBY_KMC_CFG_FILE  "/var/paas/standby_kmc_cfg.txt"


#define PAAS_AES_CRYPTO_LOG   "/var/paas/paas_crypto_log.log"
#define WSEC_AES_CRYPTO_LOG   "/var/paas/wsec_crypto_log.log"

#ifdef __cplusplus
extern "C"{
#endif /* __cplusplus */

/**  ����ֵ */
enum
{
    RET_SUCCESS         = 0,  /**< �ɹ� */
    RET_INVALID_PARAM   = 1,  /**< ��δ��� */
    RET_NORMAL_FAILURE  = 2,  /**< �ڲ�һ���쳣 */
};

/**
************************************************************
 *@ingroup
 *@brief �ӽ��ܳ�ʼ��
 *
 *
 *
 *@retval int          ��ʼ��������
  �޸���ʷ   :
  1.��    ��   : 2015��04��28��
    ��    ��   : z00223295
    �޸�����:
************************************************************/
int aesInit();

/**
************************************************************
 *@ingroup
 *@brief ע�Ṥ����Կ
 *
 *@param iDomainId       ��ID
 *@param iKeyId          keyID
 *@param pPlainTexKey    key
 *@param iKeyLen         key����
 *
 *@retval int          ע��key������
  �޸���ʷ   :
  1.��    ��   : 2015��04��28��
    ��    ��   : z00223295
    �޸�����:
************************************************************/
int registerWorkingKey(int iDomainId, int iKeyId, char *pPlainTexKey, int iKeyLen);

/**
************************************************************
 *@ingroup
 *@brief ������Կ״̬Ϊinactive����ʱ��Կֻ�������ܣ����ܼ�������
 *
 *@param iDomainId       ��ID
 *@param iKeyId          keyID
 *
 *@retval int          ����key������
  �޸���ʷ   :
  1.��    ��   : 2015��04��28��
    ��    ��   : z00223295
    �޸�����:
************************************************************/
int setKeyInvalid(int iDomainId, int iKeyId);

/**
************************************************************
 *@ingroup
 *@brief aes����
 *
 *@param iDomainId       ��ID
 *@param pcText          ����
 *
 *@retval char*          ����
  �޸���ʷ   :
  1.��    ��   : 2015��04��28��
    ��    ��   : z00223295
    �޸�����:
************************************************************/
char *aesEncrypt(int iDomainId, char *pcText);

/**
************************************************************
 *@ingroup
 *@brief aes����
 *
 *@param iDomainId       ��ID
 *@param pcText          ����
 *
 *@retval char*          ����
  �޸���ʷ   :
  1.��    ��   : 2015��04��28��
    ��    ��   : z00223295
    �޸�����:
************************************************************/
char *aesDecrypt(int iDomainId, char *pcHexEncData);

/**
************************************************************
 *@ingroup
 *@brief aes�����ļ�
 *
 *@param iDomainId       ��ID
 *@param pszPlainFile    ԭʼ�ļ�������Ҫָ��·��
 *@param pszPlainFile    �����ļ�������Ҫָ��·��
 *
 *@retval int            ���ܳɹ�ʧ�ܷ�����
  �޸���ʷ   :
  1.��    ��   : 2015��04��28��
    ��    ��   : z00223295
    �޸�����:
************************************************************/
int aesFileEncrypt(int iDomainId, const char *pszPlainFile, const char *pszCipherFile);

/**
************************************************************
 *@ingroup
 *@brief aes�����ļ�
 *
 *@param iDomainId       ��ID
 *@param pszPlainFile    �����ļ�������Ҫָ��·��
 *@param pszPlainFile    ԭʼ�ļ�������Ҫָ��·��
 *
 *@retval int            ���ܳɹ�ʧ�ܷ�����
  �޸���ʷ   :
  1.��    ��   : 2015��04��28��
    ��    ��   : z00223295
    �޸�����:
************************************************************/
int aesFileDecrypt(int iDomainId, const char *pszCipherFile, const char *pszPlainFile);

#ifdef __cplusplus
}
#endif /* __cplusplus */


#endif /* __PAAS_AES_CRYPTO_H__ */
